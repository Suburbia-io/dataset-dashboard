package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"

	"git.crumpington.com/public/am/amclient"
	"github.com/BurntSushi/toml"
	"github.com/Suburbia-io/dashboard/pkg/application"
	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/sftp"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/ssh"
)

func main() {
	var fresh bool
	flag.BoolVar(&fresh, "fresh", false, "refresh the database")
	flag.Parse()

	if fresh {
		fmt.Println("Are you sure you want to refresh the database (y/N): ")
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil || strings.ToLower(strings.TrimSpace(response)) != "y" {
			fresh = false
		}
	}

	configPath := os.Getenv("SUBURBIA_DASHBOARD_CONFIG")
	if configPath == "" {
		log.Fatalf("Environment variable SUBURBIA_DASHBOARD_CONFIG isn't set.")
	}

	config := application.Config{}

	if _, err := toml.DecodeFile(os.Getenv("SUBURBIA_DASHBOARD_CONFIG"), &config); err != nil {
		log.Fatalf("Failed to open config file: %s\n", err)
	}

	// Set error alerting.
	amClient := amclient.New(config.AlertURL)
	errors.AlertCallback = func(msg string) {
		amClient.Alert(msg)
	}

	// Periodically ping alert server.
	go func() {
		for range time.NewTicker(30 * time.Second).C {
			amClient.Ping()
		}
	}()

	services := []application.Service{
		application.DBService(fresh),
		application.MailerService,
		application.ViewService,
	}
	if config.Env == "dev" {
		services = append(services, application.DevSetupService)
	}
	app, err := application.Mount(config, services)
	if err != nil {
		errors.Unexpected.Wrap("Failed to mount application: %w", err).Alert()
		os.Exit(3)
	}

	httpSrv := &http.Server{
		ReadTimeout:  32 * time.Minute,
		WriteTimeout: 32 * time.Minute,
		Addr:         ":" + config.HTTPPort,
		Handler:      app.Routes(),
	}

	amClient.Log("Starting dashboard server...")

	go func() {
		if config.Env != "dev" {
			if config.HTTPPort != "443" {
				amClient.Alert("Not in dev, but also not using port 443.")
				os.Exit(3)
			}
			err := httpSrv.Serve(autocert.NewListener(config.Hostnames...))
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				errors.Unexpected.Wrap("HTTP server failed: %w", err).Alert()
				os.Exit(1)
			}
		} else {
			err := httpSrv.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				errors.Unexpected.Wrap("HTTP server failed: %w", err).Alert()
				os.Exit(1)
			}
		}
	}()

	// --------------------- SFTP Server ---------------------

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp4", config.SftpListenAddr)
	if err != nil {
		log.Fatal("failed to listen for connection", err)
	}

	amClient.Log("Starting SFTP server...")

	go func() {
		for {
			nConn, err := listener.Accept()
			if err != nil {
				log.Fatal("failed to accept incoming SFTP connection", err)
			}
			go handleSftp(nConn, *app.DBAL, config.SftpPrivateKey, config.SftpDir)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	<-sigChan

	if err := httpSrv.Close(); err != nil {
		errors.Unexpected.Wrap("Failed to stop server: %w", err).Alert()
	} else {
		amClient.Alert("Dashboard server stopped.")
	}
}

func handleSftp(nConn net.Conn, Dbal database.DBAL, privateKeyPath string, rootFolder string) {
	sftpServerConfig := &ssh.ServerConfig{
		PasswordCallback: nil,
		PublicKeyCallback: func(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			userName := c.User()
			user, err := Dbal.UserGetBySFTPUsername(userName)
			if err != nil {
				return nil, err
			}
			if len(user.SFTPPubKey) == 0 {
				return nil, fmt.Errorf("user has no sftp public key")
			}

			userPublicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(user.SFTPPubKey))
			if err != nil {
				println(err.Error())
				return nil, err
			}

			if bytes.Equal(userPublicKey.Marshal(), key.Marshal()) {
				// Link this SFTP session to the user logging in so we can check folder access
				// permissions for this user on this connection later on.

				sessionToken := base64.StdEncoding.EncodeToString(c.SessionID())
				sftpSession, err := Dbal.SftpSessionCreate(sessionToken, user)
				if err != nil {
					println(err.Error())
					return nil, err
				}

				// Create a login event
				defer Dbal.AuditTrailBySftpUserInsertAsync(sftpSession, tables.Customers.Table(), user.CustomerID, "SftpConnectionOpened", "")
				return nil, nil
			}

			return nil, fmt.Errorf("given public key not known for user")
		},
		BannerCallback: func(c ssh.ConnMetadata) string {
			return "\nWARNING!\nThis system is for the use of authorized users only. The activities of the users on this system are monitored and recorded by system personnel. Anyone using this system expressly consents to such monitoring and is advised that if such monitoring reveals possible evidence of criminal activity, system personnel may provide the evidence of such monitoring to law enforcement officials. Access only for persons explicitly authorized by Suburbia. All rights reserved.\n\n"
		},
		MaxAuthTries: 1,
		NoClientAuth: false,
	}

	privateBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal("Failed to load private key", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key", err)
	}

	sftpServerConfig.AddHostKey(private)

	// Before use, a handshake must be performed on the incoming
	// net.Conn.
	conn, chans, reqs, err := ssh.NewServerConn(nConn, sftpServerConfig)
	if err != nil {
		log.Printf("SFTP handshake failed: %v", err)
		return
	}

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	sessionId := base64.StdEncoding.EncodeToString(conn.SessionID())

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of an SFTP session, this is "subsystem"
		// with a payload string of "<length=4>sftp"
		log.Printf("Incoming channel: %s", newChannel.ChannelType())

		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			log.Printf("Rejected unknown channel type: %s", newChannel.ChannelType())
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("could not accept channel.", err)
		}
		log.Printf("Channel accepted")

		// Sessions have out-of-band requests such as "shell",
		// "pty-req" and "env".  Here we handle only the
		// "subsystem" request.
		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok := false
				switch req.Type {
				case "subsystem":
					if string(req.Payload[4:]) == "sftp" {
						ok = true
					}
				}
				req.Reply(ok, nil)
			}
		}(requests)

		serverOptions := []sftp.ServerOption{
			sftp.ReadOnly(),
		}

		server, err := sftp.NewServer(
			channel,
			rootFolder,
			sessionId,
			Dbal,
			serverOptions...,
		)

		if err != nil {
			log.Fatal(err)
		}
		if err := server.Serve(); err == io.EOF {
			server.Close()
			log.Print("sftp client exited session.")
		} else if err != nil {
			log.Fatal("sftp server completed with error:", err)
		}
	}
}
