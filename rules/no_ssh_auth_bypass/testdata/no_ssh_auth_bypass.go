package fixtures

import (
	"fmt"
	"ssh"
)

// Invalid: PublicKeyCallback sets Extensions in ssh.Permissions (grants roles based on unverified key)
func badSetupServer() *ssh.ServerConfig {
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) { // MATCH /PublicKeyCallback sets Extensions in ssh.Permissions, which may lead to authentication bypass/
			if isAuthorizedKey(key) {
				return &ssh.Permissions{
					Extensions: map[string]string{
						"role": "admin",
					},
				}, nil
			}
			return nil, fmt.Errorf("unauthorized")
		},
	}
	return config
}

// Valid: PublicKeyCallback returns empty ssh.Permissions without Extensions
func goodSetupServer() *ssh.ServerConfig {
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if isAuthorizedKey(key) {
				return &ssh.Permissions{}, nil
			}
			return nil, fmt.Errorf("unauthorized")
		},
	}
	return config
}

// Valid: PublicKeyCallback returns nil Permissions
func goodNilPermissions() *ssh.ServerConfig {
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if isAuthorizedKey(key) {
				return nil, nil
			}
			return nil, fmt.Errorf("unauthorized")
		},
	}
	return config
}

// Valid: No PublicKeyCallback field set
func goodNoCallback() *ssh.ServerConfig {
	config := &ssh.ServerConfig{}
	return config
}

// Valid: PublicKeyCallback with Extensions as empty map
func goodEmptyExtensions() *ssh.ServerConfig {
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if isAuthorizedKey(key) {
				return &ssh.Permissions{
					Extensions: map[string]string{},
				}, nil
			}
			return nil, fmt.Errorf("unauthorized")
		},
	}
	return config
}

func isAuthorizedKey(key ssh.PublicKey) bool {
	return true
}
