# Intel Project Amber Go SGX Adapter
This is the beta version of Go SGX Adapter for collecting Quote from SGX enabled platform.

This library leverages Intel SGX DCAP for Quote generation: [https://github.com/intel/SGXDataCenterAttestationPrimitives](https://github.com/intel/SGXDataCenterAttestationPrimitives)

## Go Requirement

Use <b>go1.19 or newer</b>. Follow https://go.dev/doc/install for installation of Go.

## Usage

Create a new Go SGX adapter, then use the adapter to
collect quote from SGX enabled platform.

```go
import "github.com/intel/amber-client/go-sgx"

adapter, err := sgx.NewAdapter(enclaveId, enclaveHeldData, unsafe.Pointer(C.enclave_create_report))

evidence, err := adapter.CollectEvidence(nonce)
if err != nil {
    return err
}
```

## License

This library is distributed under the BSD-style license found in the [LICENSE](../LICENSE)
file.
