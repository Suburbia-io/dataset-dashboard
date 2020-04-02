import {adminRpcClient} from '../../shared/js/rpcClient'

let client = null

export default {
  init: (baseURL) => client = adminRpcClient(baseURL),
  datasetTagTypeList: ({datasetID}) => client('fingerprintListTags', {
    datasetID
  }),
  fingerprintList: ({datasetID, tagAppID, limit, offset, orderBy, orderAsc, countThreshold, fingerprintIncludes, fingerprintExcludes, consTagIncludes, consTagExcludes, rawTextIncludes, rawTextExcludes, tagIncludes, tagExcludes,}) => client('fingerprintList', {
    datasetID,
    tagAppID,
    limit,
    offset,
    orderBy,
    orderAsc,
    countThreshold,
    fingerprintIncludes,
    fingerprintExcludes,
    consTagIncludes,
    consTagExcludes,
    rawTextIncludes,
    rawTextExcludes,
    tagIncludes,
    tagExcludes,
  }),
  fingerprintTagSuggestions: ({datasetID, tagTypeID, search}) => client('fingerprintTagSuggestions', {
    datasetID,
    tagTypeID,
    search,
  }),
  fingerprintUpsertTags: ({datasetID, tagTypeID, tagValue, fingerprints}) => client('fingerprintUpsertTags', {
    datasetID,
    tagTypeID,
    tagValue,
    fingerprints,
  }),
  fingerprintsListTagApps: () => client('fingerprintsListTagApps')
}
