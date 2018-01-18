package conf

import (
  "fmt"
  "github.com/spf13/viper"
)

// Read all needed configs if no errors occurred return them as map, do not change the order unless you 100% sure
func ReadConfig() map[string]interface{} {
  addDefaults()
  readConfJson()
  readEnvVariables()
  return viper.AllSettings()
}

// Determine path to work dir and add it to Viper.
// Path is determined as $GOPATH/src/go-team-room
func setupWorkDir() {
  readVar("go_path", "GOPATH")
  fmt.Println("GOPATH is: " + viper.GetString("work_dir"))
  viper.SetDefault("work_dir", viper.GetString("go_path")+"/src/go-team-room")
  fmt.Println("WorkDir is: " + viper.GetString("work_dir"))
}

// Add default when possible
func addDefaults() {
  setupWorkDir()
  viper.SetDefault("ip", "127.0.0.1")
  viper.SetDefault("port", 8080)
  viper.SetDefault("static_dir", viper.GetString("workDir") + "/client/dist")
}

//Read configs from conf json, if cant read error occurred.
// Path for read is $GOPATH/src/go-team-room/conf/conf.json
func readConfJson() {
  viper.SetConfigName("conf")
  viper.SetConfigType("json")
  viper.AddConfigPath(viper.GetString("work_dir") + "/conf")
  fmt.Println("Start reading config from conf.json")
  err := viper.ReadInConfig()
  if err != nil { // Handle errors reading the config file
    panic(fmt.Errorf("Cant read config file, error is: %s \n", err))
  }
  fmt.Print("Finished reading conf.json, current configs: ")
  fmt.Println(viper.AllSettings())
}

// Read environment variables into viper, if they not present error occurred.
func readEnvVariables() {
  readVar("aws_access_key_id", "GO_AWS_ACCESS_KEY_ID")
  readVar("aws_secret_key", "GO_AWS_SECRET_ACCESS_KEY")
  readVar("dynamo_endpoint", "GO_AWS_DYNAMO_ENDPOINT")
  readVar("dynamo_region", "GO_AWS_DYNAMO_REGION")
  readVar("mysql_dsn", "GO_MYSQL_DSN")
}

func readVar(key, name string) {
  fmt.Printf("Read %s from environment \n", name)
  err := viper.BindEnv(key, name)
  if err != nil {
    panic(fmt.Errorf("Error ocured while reading %s, error is: %s \n", name, err))
  } else if viper.GetString(key) == "" {
    panic(fmt.Errorf("%s is missing in your system, you must setup it", name))
  }
}
