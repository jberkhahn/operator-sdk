package generate

type PackageManifestGenerator interface {
	Generate(opts *PkgOptions) error
}
