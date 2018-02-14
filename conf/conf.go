package conf

var allConfigs = ReadConfig()

var LoginUrl = allConfigs["login_url"].(string)
var WorkDir = allConfigs["work_dir"].(string)
var Ip = allConfigs["ip"].(string)
var Port = allConfigs["port"].(string)
var StaticDir = allConfigs["static_dir"].(string)
var AwsAccessKeyId = allConfigs["aws_access_key_id"].(string)
var AwsSecretKey = allConfigs["aws_secret_key"].(string)
var AwsBucketName = allConfigs["aws_bucket_name"].(string)
var DynamoEndpoint = allConfigs["dynamo_endpoint"].(string)
var DynamoRegion = allConfigs["dynamo_region"].(string)
var MysqlDsn = allConfigs["mysql_dsn"].(string)
var MysqlDBName = allConfigs["my_sql_db_name"].(string)
var MaxMessages = int(allConfigs["max_messages_at_once"].(float64))
var MaxChatRooms = int(allConfigs["max_chatrooms"].(float64))

var EnableSendMails = allConfigs["enable_email"].(bool)
var GohumEmail = allConfigs["gohum_email"].(string)
var GohumEmailPass = allConfigs["gohum_email_password"].(string)
var SmtpServer = allConfigs["smtp_server"].(string)
var SmtpPort = int(allConfigs["smtp_port"].(float64))
