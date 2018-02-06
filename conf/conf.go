package conf

var allConfigs = ReadConfig()

var WorkDir = allConfigs["work_dir"].(string)
var Ip = allConfigs["ip"].(string)
var Port = allConfigs["port"].(string)
var StaticDir = allConfigs["static_dir"].(string)
var AwsAccessKeyId = allConfigs["aws_access_key_id"].(string)
var AwsSecretKey = allConfigs["aws_secret_key"].(string)
var DynamoEndpoint = allConfigs["dynamo_endpoint"].(string)
var DynamoRegion = allConfigs["dynamo_region"].(string)
var MysqlDsn = allConfigs["mysql_dsn"].(string)
var MaxMessages = int(allConfigs["max_messages_at_once"].(float64))
var MaxChatRooms = int(allConfigs["max_chatrooms"].(float64))
