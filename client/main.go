package main

import (
	"context"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/credentials"
	"log"
	"time"

	pb "github.com/aau-network-security/haaukins-store/proto"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	AUTH_KEY    = "au"
)

var (
	UnreachableDaemonErr = errors.New("Daemon seems to be unreachable")
	UnauthorizedErr      = errors.New("You seem to not be logged in")
)

type Creds struct {
	Token    string
	Insecure bool
}

func (c Creds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"token": string(c.Token),
	}, nil
}

func (c Creds) RequireTransportSecurity() bool {
	return !c.Insecure
}

func main(){

	//todo just for test purpose
	test_auth_key := "c41ec030-db76-473f-a504-5a7323aa04ec"
	test_sign_key := "34b16c10-1a2c-4533-83e8-cfde78817501"
	test_cert_path := "/home/gian/Documents/haaukins_files/server.crt"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		AUTH_KEY: test_auth_key,
	})

	tokenString, err := token.SignedString([]byte(test_sign_key))
	if err != nil {
		fmt.Println("Error creating the token")
	}


	authCreds := Creds{Token: tokenString}
	dialOpts := []grpc.DialOption{}

	ssl := true
	if ssl {
		creds, _ := credentials.NewClientTLSFromFile(test_cert_path, "")
		dialOpts = append(dialOpts,
			grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(authCreds))
	} else {
		authCreds.Insecure = true
		dialOpts = append(dialOpts,
			grpc.WithInsecure(),
			grpc.WithPerRPCCredentials(authCreds))
	}

	conn, err := grpc.Dial(address, dialOpts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewStoreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Println(tokenString)

	//addEventRequest := pb.AddEventRequest{
	//	Name:                 "Test from Client",
	//	Tag:                  "clientestttttt",
	//	Frontends:            "awdwad,wadwad,rtr,trt",
	//	Exercises:            "bla,bla1,ciao",
	//	Available:            1212,
	//	Capacity:             20,
	//	ExpectedFinishTime:   "wadwad wdawadwadwa  awdadwad adwd",
	//}

	//addTeam := pb.AddTeamRequest{
	//	Id:                   "its_working",
	//	EventTag:             "menne",
	//	Email:                "menne@menne.com",
	//	Name:                 "menne",
	//	Password:             "menne_token_test",
	//}
	//r, err := c.AddTeam(ctx, &addTeam)

	//r, err := c.UpdateTeamSolvedChallenge(ctx, &pb.UpdateTeamSolvedChallengeRequest{
	//	TeamId:               "menne2",
	//	Tag:                  "prova",
	//	CompletedAt:          "prova time",
	//})
	r, err := c.GetEvents(ctx, &pb.EmptyRequest{})
	if err != nil{
		log.Fatalf("could not greet: %v", err)
	}
	if r.ErrorMessage != ""{
		log.Fatalf("my could not greet: %v", r.ErrorMessage)
	}
	//log.Println(r.Message)
	for _, e := range r.Events{
		fmt.Println(e)
	}
}