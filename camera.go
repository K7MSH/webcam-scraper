package main

type Camera struct {
	Name   string
	URL    string
	SaveTo string
	Auth   *CameraAuth
}
type CameraAuth struct {
	CameraName string
	User       string
	Password   string
}
