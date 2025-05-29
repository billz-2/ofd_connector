# billz_ofd_connector

package for calling ofd endpoints

- how to add new package version:

  - push changes into master
  - git tag -a v0.0.13 -m "skip stack level"
  - git push origin v0.0.13
  - wait some time (up to 30 minutes) for proxy to update
  - to update package in microservice execute `go get -u github.com/billz-2/ofd_connector`

- how to execute tests
  - go test ./...
