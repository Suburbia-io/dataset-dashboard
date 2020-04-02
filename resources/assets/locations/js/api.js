import {adminRpcClient} from '../../shared/js/rpcClient'
let client = null;

export default {
  init: (baseURL) => client = adminRpcClient(baseURL),
  locationApprove: ({datasetID, locationHash, approved}) => client('location-approve/', {
    datasetID,
    locationHash,
    approved,
  }),
  locationSetGeonameID: ({datasetID, locationHash, geonameID}) => client('location-set-geoname-id/', {
    datasetID,
    locationHash,
    geonameID,
  }),
}
