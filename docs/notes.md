# Open questions
* How do we use Podman bindings with Windows?
  * (A) [Use `podman info --debug` to find information about where podman socket is on WSL](https://github.com/containers/podman/issues/13246)
  * How do you determine if a go binary is running on WSL?
  * `podman info --debug` shows that the podman socket is in a location which I can't find on WSL - why is that?
* How can we test running `beacon` on each environment?

# References
* [Guide on setting up Podman bindings for Linux (through Go)](https://podman.io/blogs/2020/08/10/podman-go-bindings.html#connect-service)
* [Find out which platform go is running in during runtime using GOOS](https://pkg.go.dev/runtime#GOOS)