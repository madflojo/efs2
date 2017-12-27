package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"testing"
)

func Test_initSSH_NoKeyfile(t *testing.T) {
	var opt options
	opt.Keyfile = "/dev/null"
	_, err := initSSH(opt)
	if err == nil {
		t.Errorf("initSSH should have returned an error but did not")
	}
}

func Test_initSSH_UsernameOpt(t *testing.T) {
	var opt options
	opt.Keyfile = "/tmp/key"
	b := []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA0iCf8/aVYKonGlBfx0ceQ3wQemRISwyiMX5YZ8iL5ZN5vj6N\naVHK9dtgMHJCLfXaZocdneaBrQ6La9Lz32dWgCcKGjchN2fh+0QN6EkkY1P/ijvO\n3YLNIOawnD2HlpzhfEUzhdayshaMwHm6cxig7KcJMCwT5ntlcymcqQCZqTx5+vpO\nUPmY9kP6qLSOn5PUX49daDD3CH3kdu+KOxk8L7gukdo9RJ2EIkDH+m5xVnMarGGt\nv/bbbMT6aEEpmeTzA+MmBt62tDGcPmsFINacGTKbKJ9EClGUyvU/l8fxZLgDFots\nccXe8NOALGRSaohJqW3JC35kOnUlYa3Jt2yTSwIDAQABAoIBAQC0O4GXS1kDSc8y\ndeBBWJHvxnmH0X5kyRhRpZKEqnK8Xwucj6DRxnN1AE74Hvj+3RMQwDI6Ht35pzEV\nMiM16zg5wcKbi8/06yjdUZkwNZR9ki3szrH4M9porxarXOdw221ZHy47TVWHBWqD\nKaYwVN6rPfbWl+gV2J/C8N1L5JTon8YMOUHu4MftfihA5S8dUU0eUqbKrag7MjZf\nplXq69jtIpTBHPqP5kTyzanDams3HiocAYczEJ1h9ZshrSH2VTjnP8DOFOrB0Tro\nUfgS9w/M72UX4YBY+TEclEA3lyXNY/LZx/7nRB9NSrRKp+ddmrn6uTf9Sii3EoZ+\nX/H7x4uZAoGBAOiIYbcKADRv8BNQFP590nOxVWXfh6E3xvNX6wknxGTCQ+toFm8j\nQmQRNd50qxZjTtX5bMFU+T/opB6dNJiPpCvSJxhgfrzpAmac2bDhM7uwqRBIhUq1\nTpW3x/+nxPB91QLd03w/VrDMeF9DweoQ8k7hy73KQ9h92WhFWIujp86PAoGBAOdV\nY6Fo1dSeRdnTJkhVhkpq0U3OKgjg/lfW9Gl2F8RfjX4bvCDiLqHgYj5kee+GakqQ\nX98H3L0Ln0e4OvX9IulEFas1hoEVcG4b0Gt+N9jYD9j2oem6QGh7Eo8e2oZbiD2w\nh1VvAPkLNnTLMh3t0SqrM2c7VcLbzYozZcSIrw2FAoGAZjQTZsURg/qzNXncUGLz\nDgCifU07KsP+QNSirHp7GqI8AgrU2XJQ6vSZjbYPwJ5Tdz7S60Ky7sEM6ZvFE00H\nJm+O2WsIKXXspTdJgzHocBVcqZiGZWi9KpcFY7vUlrNn3YOsQY8BRmIIgi6g24Up\nSzx5NWjiWxQta2QXYADFb8cCgYEAx4QK0KxFOAJlhi+pJdu1Xbtw12UHNe8vDf1T\ngR2b8/7hXrF2+Pl6dJy8vskTrXTFeZe5R/dU9yrt7gJDv+LZ2EujUK1yWyRtelm8\n0OjkK751NI/KJ3Y+lJ4I7K0Ulaqd/26f2hxJv4FfLy4NBGkW7HEJfXBcUBoY1Kft\njqrAJ60CgYBnVlikLRYW5WcHUKD+7OiEzU6KoT27p/i3LVr5RQ+s1mMrV5vXJt3b\nPwGGq5vVJ/fi5gt9ezAWOaL5l11ALI3tnlIEAaRIsPF+hX6XbisL7DqPJtqwu095\nI3IKMZT1Aom6EaY+HMf5Y72//WDTROD/h69NAlOoKqR4lpw69icEcA==\n-----END RSA PRIVATE KEY-----\n")
	err := ioutil.WriteFile(opt.Keyfile, b, 0644)
	if err != nil {
		t.Errorf("Error creating test keyfile")
	}
	opt.User = "testuser"
	s, err := initSSH(opt)
	_ = os.Remove(opt.Keyfile)
	if err != nil {
		t.Errorf("initSSH returned an unexpected error - %s", err)
	}
	if s.User != "testuser" {
		t.Errorf("initSSH returned an invalid user - %s", s.User)
	}
}

func Test_initSSH(t *testing.T) {
	var opt options
	opt.Keyfile = "/tmp/key"
	b := []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA0iCf8/aVYKonGlBfx0ceQ3wQemRISwyiMX5YZ8iL5ZN5vj6N\naVHK9dtgMHJCLfXaZocdneaBrQ6La9Lz32dWgCcKGjchN2fh+0QN6EkkY1P/ijvO\n3YLNIOawnD2HlpzhfEUzhdayshaMwHm6cxig7KcJMCwT5ntlcymcqQCZqTx5+vpO\nUPmY9kP6qLSOn5PUX49daDD3CH3kdu+KOxk8L7gukdo9RJ2EIkDH+m5xVnMarGGt\nv/bbbMT6aEEpmeTzA+MmBt62tDGcPmsFINacGTKbKJ9EClGUyvU/l8fxZLgDFots\nccXe8NOALGRSaohJqW3JC35kOnUlYa3Jt2yTSwIDAQABAoIBAQC0O4GXS1kDSc8y\ndeBBWJHvxnmH0X5kyRhRpZKEqnK8Xwucj6DRxnN1AE74Hvj+3RMQwDI6Ht35pzEV\nMiM16zg5wcKbi8/06yjdUZkwNZR9ki3szrH4M9porxarXOdw221ZHy47TVWHBWqD\nKaYwVN6rPfbWl+gV2J/C8N1L5JTon8YMOUHu4MftfihA5S8dUU0eUqbKrag7MjZf\nplXq69jtIpTBHPqP5kTyzanDams3HiocAYczEJ1h9ZshrSH2VTjnP8DOFOrB0Tro\nUfgS9w/M72UX4YBY+TEclEA3lyXNY/LZx/7nRB9NSrRKp+ddmrn6uTf9Sii3EoZ+\nX/H7x4uZAoGBAOiIYbcKADRv8BNQFP590nOxVWXfh6E3xvNX6wknxGTCQ+toFm8j\nQmQRNd50qxZjTtX5bMFU+T/opB6dNJiPpCvSJxhgfrzpAmac2bDhM7uwqRBIhUq1\nTpW3x/+nxPB91QLd03w/VrDMeF9DweoQ8k7hy73KQ9h92WhFWIujp86PAoGBAOdV\nY6Fo1dSeRdnTJkhVhkpq0U3OKgjg/lfW9Gl2F8RfjX4bvCDiLqHgYj5kee+GakqQ\nX98H3L0Ln0e4OvX9IulEFas1hoEVcG4b0Gt+N9jYD9j2oem6QGh7Eo8e2oZbiD2w\nh1VvAPkLNnTLMh3t0SqrM2c7VcLbzYozZcSIrw2FAoGAZjQTZsURg/qzNXncUGLz\nDgCifU07KsP+QNSirHp7GqI8AgrU2XJQ6vSZjbYPwJ5Tdz7S60Ky7sEM6ZvFE00H\nJm+O2WsIKXXspTdJgzHocBVcqZiGZWi9KpcFY7vUlrNn3YOsQY8BRmIIgi6g24Up\nSzx5NWjiWxQta2QXYADFb8cCgYEAx4QK0KxFOAJlhi+pJdu1Xbtw12UHNe8vDf1T\ngR2b8/7hXrF2+Pl6dJy8vskTrXTFeZe5R/dU9yrt7gJDv+LZ2EujUK1yWyRtelm8\n0OjkK751NI/KJ3Y+lJ4I7K0Ulaqd/26f2hxJv4FfLy4NBGkW7HEJfXBcUBoY1Kft\njqrAJ60CgYBnVlikLRYW5WcHUKD+7OiEzU6KoT27p/i3LVr5RQ+s1mMrV5vXJt3b\nPwGGq5vVJ/fi5gt9ezAWOaL5l11ALI3tnlIEAaRIsPF+hX6XbisL7DqPJtqwu095\nI3IKMZT1Aom6EaY+HMf5Y72//WDTROD/h69NAlOoKqR4lpw69icEcA==\n-----END RSA PRIVATE KEY-----\n")
	err := ioutil.WriteFile(opt.Keyfile, b, 0644)
	if err != nil {
		t.Errorf("Error creating test keyfile")
	}
	s, err := initSSH(opt)
	_ = os.Remove(opt.Keyfile)
	if err != nil {
		t.Errorf("initSSH returned an unexpected error - %s", err)
	}

	u, _ := user.Current()
	if s.User != u.Username {
		t.Errorf("initSSH returned an unexpected user - %s", err)
	}
}

func Test_readKeyfile_NoFile(t *testing.T) {
	_, err := readKeyfile("/nothing")
	if err == nil {
		t.Errorf("readKeyfile did not return the expected error, returned nil")
	}
}

func Test_readKeyfile(t *testing.T) {
	f := "/tmp/key"
	b := []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA0iCf8/aVYKonGlBfx0ceQ3wQemRISwyiMX5YZ8iL5ZN5vj6N\naVHK9dtgMHJCLfXaZocdneaBrQ6La9Lz32dWgCcKGjchN2fh+0QN6EkkY1P/ijvO\n3YLNIOawnD2HlpzhfEUzhdayshaMwHm6cxig7KcJMCwT5ntlcymcqQCZqTx5+vpO\nUPmY9kP6qLSOn5PUX49daDD3CH3kdu+KOxk8L7gukdo9RJ2EIkDH+m5xVnMarGGt\nv/bbbMT6aEEpmeTzA+MmBt62tDGcPmsFINacGTKbKJ9EClGUyvU/l8fxZLgDFots\nccXe8NOALGRSaohJqW3JC35kOnUlYa3Jt2yTSwIDAQABAoIBAQC0O4GXS1kDSc8y\ndeBBWJHvxnmH0X5kyRhRpZKEqnK8Xwucj6DRxnN1AE74Hvj+3RMQwDI6Ht35pzEV\nMiM16zg5wcKbi8/06yjdUZkwNZR9ki3szrH4M9porxarXOdw221ZHy47TVWHBWqD\nKaYwVN6rPfbWl+gV2J/C8N1L5JTon8YMOUHu4MftfihA5S8dUU0eUqbKrag7MjZf\nplXq69jtIpTBHPqP5kTyzanDams3HiocAYczEJ1h9ZshrSH2VTjnP8DOFOrB0Tro\nUfgS9w/M72UX4YBY+TEclEA3lyXNY/LZx/7nRB9NSrRKp+ddmrn6uTf9Sii3EoZ+\nX/H7x4uZAoGBAOiIYbcKADRv8BNQFP590nOxVWXfh6E3xvNX6wknxGTCQ+toFm8j\nQmQRNd50qxZjTtX5bMFU+T/opB6dNJiPpCvSJxhgfrzpAmac2bDhM7uwqRBIhUq1\nTpW3x/+nxPB91QLd03w/VrDMeF9DweoQ8k7hy73KQ9h92WhFWIujp86PAoGBAOdV\nY6Fo1dSeRdnTJkhVhkpq0U3OKgjg/lfW9Gl2F8RfjX4bvCDiLqHgYj5kee+GakqQ\nX98H3L0Ln0e4OvX9IulEFas1hoEVcG4b0Gt+N9jYD9j2oem6QGh7Eo8e2oZbiD2w\nh1VvAPkLNnTLMh3t0SqrM2c7VcLbzYozZcSIrw2FAoGAZjQTZsURg/qzNXncUGLz\nDgCifU07KsP+QNSirHp7GqI8AgrU2XJQ6vSZjbYPwJ5Tdz7S60Ky7sEM6ZvFE00H\nJm+O2WsIKXXspTdJgzHocBVcqZiGZWi9KpcFY7vUlrNn3YOsQY8BRmIIgi6g24Up\nSzx5NWjiWxQta2QXYADFb8cCgYEAx4QK0KxFOAJlhi+pJdu1Xbtw12UHNe8vDf1T\ngR2b8/7hXrF2+Pl6dJy8vskTrXTFeZe5R/dU9yrt7gJDv+LZ2EujUK1yWyRtelm8\n0OjkK751NI/KJ3Y+lJ4I7K0Ulaqd/26f2hxJv4FfLy4NBGkW7HEJfXBcUBoY1Kft\njqrAJ60CgYBnVlikLRYW5WcHUKD+7OiEzU6KoT27p/i3LVr5RQ+s1mMrV5vXJt3b\nPwGGq5vVJ/fi5gt9ezAWOaL5l11ALI3tnlIEAaRIsPF+hX6XbisL7DqPJtqwu095\nI3IKMZT1Aom6EaY+HMf5Y72//WDTROD/h69NAlOoKqR4lpw69icEcA==\n-----END RSA PRIVATE KEY-----\n")
	err := ioutil.WriteFile(f, b, 0644)
	if err != nil {
		t.Errorf("Error creating test keyfile")
	}
	_, err = readKeyfile(f)
	_ = os.Remove(f)
	if err != nil {
		t.Errorf("readKeyfile returned an unexpected error - %s", err)
	}
}
