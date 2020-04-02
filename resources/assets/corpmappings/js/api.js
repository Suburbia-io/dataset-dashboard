import {adminRpcClient} from '../../shared/js/rpcClient'

let client = null

export default {
  init: (baseURL) => client = adminRpcClient(baseURL),
  listTagsForTagType: ({DatasetID, TagTypeID}) => client('listTagsForTagType', {DatasetID, TagTypeID}),
  insertCorpMapping: ({corpTypeID, tagTypeID, tagID}) => client('insertCorpMapping', {corpTypeID, tagTypeID, tagID}),
  insertCorpMappingRule: ({corpMappingID, corpID, fromDate, country}) => client('insertCorpMappingRule', {
    corpMappingID,
    corpID,
    fromDate,
    country
  }),
  deleteCorpMappingRule: ({corpMappingRuleID}) => client('deleteCorpMappingRule', {corpMappingRuleID}),
}
