//
// 3nigm4 storageservice package
// Author: Guido Ronchetti <dyst0ni3@gmail.com>
// v1.0 16/06/2016
//
package main

// Go standard libraries
import (
	"fmt"
	"net/rpc"
)

// 3n4 libraries
import (
	"github.com/nexocrew/3nigm4/lib/auth"
)

// AuthClient is the interface used to interact
// with authentication services.
type AuthClient interface {
	Login(string, string) ([]byte, error)                          // manage user's login;
	Logout([]byte) ([]byte, error)                                 // manage user's logout;
	AuthoriseAndGetInfo([]byte) (*auth.UserInfoResponseArg, error) // returns authenticated user infos or an error;
	Close() error                                                  // closes eventual connections.
}

// AuthRPC implements the RPC default client for
// the 3nigm4 auth service.
type AuthRPC struct {
	client *rpc.Client
}

// NewAuthRPC creates a new instance of the RPC
// client used to interact with the auth service.
func NewAuthRPC(addr string, port int) (*AuthRPC, error) {
	address := fmt.Sprintf("%s:%d", addr, port)
	rawClient, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}
	return &AuthRPC{
		client: rawClient,
	}, nil
}

// Login grant access to users, over RPC, using username and password.
func (a *AuthRPC) Login(username string, password string) ([]byte, error) {
	// perform login on RPC service
	var loginResponse auth.LoginResponseArg
	err := a.client.Call("Login.Login", &auth.LoginRequestArg{
		Username: username,
		Password: password,
	}, &loginResponse)
	if err != nil {
		return nil, err
	}
	return loginResponse.Token, nil
}

// Logout remove actual active sessions over RPC.
func (a *AuthRPC) Logout(token []byte) ([]byte, error) {
	var logoutResponse auth.LogoutResponseArg
	err := a.client.Call("Login.Logout", &auth.LogoutRequestArg{
		Token: token,
	}, &logoutResponse)
	if err != nil {
		return nil, err
	}
	return logoutResponse.Invalidated, nil
}

// AuthoriseAndGetInfo if the token is valid returns info about
// the associated user over RPC service.
func (a *AuthRPC) AuthoriseAndGetInfo(token []byte) (*auth.UserInfoResponseArg, error) {
	// verify token and retrieve user infos
	var authResponse auth.UserInfoResponseArg
	err := a.client.Call("SessionAuth.UserInfo", &auth.AuthenticateRequestArg{
		Token: token,
	}, &authResponse)
	if err != nil {
		return nil, err
	}
	return &authResponse, nil
}

// Close closes RPC connection.
func (a *AuthRPC) Close() error {
	return a.client.Close()
}
