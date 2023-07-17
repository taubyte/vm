# Taubyte WebAssembly Virtual Machine (TVM)

The Taubyte WebAssembly Virtual Machine (TVM) serves as a crucial component in the execution layer of any Taubyte-based Cloud Computing Network. In addition, it is commonly utilized in testing scenarios, particularly when building plugins, or as we like to call them, satellites. For more details, check out [Orbit](https://github.com/taubyte/vm-orbit).

## Installation

You can easily install TVM by using the following command:

```bash
go get github.com/taubyte/vm
```

## Structure

TVM is composed of several components:

  - `backend/`: Houses backends that implement the `vm.Backend` interface. For more information, refer to [interfaces](https://github.com/taubyte/go-interfaces/vm).
  - `context/`: Implements the `vm.Context` interface, which defines the execution context of a WebAssembly module.
  - `resolvers/`: Contains resolvers that implement the `vm.Resolver`. The resolver is utilized to translate a module name into a [multiaddress](https://github.com/multiformats/multiaddr).
  - `loader/`: Implements the `vm.Loader` interface, which combines the resolver and various backends. Once a module name is resolved within a given context, the loader will loop over available backends until the module is found.
  - `service/`: Implements a WebAssembly service capable of provisioning runtimes. It adheres to the `vm.Service` specification.

## Plugins

To extend the capabilities of the Taubyte VM, you can use [Orbit](https://github.com/taubyte/vm-orbit). If you require more direct access to the Taubyte node (odo), you can derive inspiration from the [Core Plugins](https://github.com/taubyte/vm-core-plugins).

## License

This project is licensed under the BSD 3-Clause License. For more details, please refer to the LICENSE file.

## Help

Join our [Discord](https://discord.gg/taubyte) community if you need assistance or want to engage with us!

## Maintainers

 - Samy Fodil (@samyfodil)
 - Tafseer Khan (@tafseer-khan)