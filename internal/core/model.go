package core

type Configurations struct {
	KafkaConfig    	KafkaConfigurations
}

type KafkaConfigurations struct {
    Username 		string 
    Password 		string 
    Protocol        string
    Mechanisms      string
    Clientid 		string 
    Brokers 		string 
    Brokers1 		string 
    Brokers2 		string 
    Brokers3 		string 
    Groupid 		string 
	Topic  			string
    Partition       int
    ReplicationFactor int
    RequiredAcks    int
    Lag             int
}