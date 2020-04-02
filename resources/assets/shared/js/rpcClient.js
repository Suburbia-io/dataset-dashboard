export function rpcClient(host) {
    host = host.replace(/\/+$/, '')

    return function(endpoint, data = {}) {
        endpoint = endpoint.replace(/^\/+/, '')

        return new Promise((resolve, reject) => {
            fetch(host + '/' + endpoint, {
                credentials: 'include',
                mode: 'cors',
                headers: { 'Content-Type': 'application/json' },
                method: 'POST',
                body: JSON.stringify(data)
            })
                .catch(() => reject('Unexpected'))
                .then((response) => {
                    if (!response || !response.ok) return reject('Unexpected')

                    response.json()
                        .then(({ ok, data }) => {
                            return ok ? resolve(data) : reject(data)
                        })
                        .catch(() => reject('Unexpected'))
                })
        })
    }
}

export function adminRpcClient(baseURL) {
  const client = rpcClient(baseURL);
  return function (endpoint, data = {}) {
    return client(endpoint, data).catch(err => {
      if (err === 'HttpNotAuthorized') {
        window.location = '/admin/auth/'
      }
      throw err
    })
  }
}
