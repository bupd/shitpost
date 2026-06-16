# Release security plan

This plan turns OWASP release and supply-chain guidance into concrete work for `shitpost`.

The goal is not to claim security by adding badges. The goal is to make every release traceable, reviewable, reproducible enough to verify, and difficult to tamper with between source and published artifact.

## OWASP baseline

The relevant OWASP guidance comes from:

- [OWASP CI/CD Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/CI_CD_Security_Cheat_Sheet.html)
- [OWASP Software Supply Chain Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Software_Supply_Chain_Security_Cheat_Sheet.html)
- [OWASP Software Component Verification Standard](https://owasp.org/www-project-software-component-verification-standard/)
- [OWASP SAMM](https://owasp.org/www-project-samm/)

For this repo, those sources reduce to these release requirements:

- Protect source control with reviewed pull requests, protected branches, protected release tags, MFA, and no bypassable direct production release path.
- Run builds only on hosted, isolated, ephemeral CI workers for release artifacts.
- Grant CI jobs the least permissions they need, and separate build, signing, publishing, and deployment authority where practical.
- Keep secrets out of source, logs, images, binaries, and release artifacts.
- Pin dependencies and third-party CI actions; review and update them through an explicit process.
- Generate SBOMs and monitor dependencies, base images, build tools, and release artifacts for vulnerabilities.
- Sign artifacts and container images; publish checksums and provenance bound to artifact digests.
- Verify provenance and signatures before consumers deploy or mirror release artifacts.
- Keep release logs and audit events useful for incident response without leaking credentials.
- Document a release runbook and recovery path for compromised credentials, compromised dependencies, or bad releases.

## Current state

Already present:

- GitHub Actions runs Go formatting, vet, build, and tests on pull requests through `task ci`.
- Release archives are built by GoReleaser on version tags.
- GoReleaser emits checksums and SBOMs for archives.
- Docker Buildx emits SBOM and provenance metadata for container images.
- GHCR images are signed with keyless `cosign`.
- SLSA generic provenance is configured for release archives.
- SLSA container provenance is configured for GHCR images.
- `docs/slsa.md` documents the SLSA Build L3 target and verification commands.

Known gaps:

- GitHub branch and tag protection are not documented as required release controls.
- Third-party GitHub Actions are version-tagged but not pinned to immutable SHAs, except the SLSA reusable workflows which must remain semver-tagged for verifier compatibility.
- No repo-local Dependabot configuration is present for Go modules, GitHub Actions, or Docker base images.
- No explicit secret scanning gate is present in CI.
- No vulnerability gate runs `govulncheck`, OSV, Trivy, Grype, or equivalent scanners.
- No release smoke job verifies published signatures and SLSA provenance after release.
- No incident response checklist exists for bad releases or leaked publishing credentials.

## Target release model

Release artifacts:

- GitHub release archives produced by GoReleaser.
- GHCR container image published by digest and tag.

Release authority:

- Source changes land through reviewed pull requests to `main`.
- Version tags are protected and created only from reviewed commits on `main`.
- CI uses minimal `GITHUB_TOKEN` permissions per job.
- Releases are triggered only by pushed `v*` tags.

Consumer verification:

- Archives are verified with `slsa-verifier verify-artifact` and GoReleaser checksums.
- GHCR images are verified with `slsa-verifier verify-image` and `cosign verify`.

Pull request test model:

- `.github/workflows/ci.yml` runs on pull requests to `main`.
- The test job runs in the `golang:1.25-bookworm` container.
- The job installs `task` and runs `task ci`.
- `task ci` checks formatting, vets code, runs tests, and builds the bot without modifying files.

## Implementation plan

### Phase 1: lock the release path

Outcome: only reviewed source can become a release.

- Enable GitHub branch protection for `main`.
- Require pull request review before merge.
- Require status checks: formatting, vet, build, tests, workflow lint.
- Disable direct pushes to `main` except maintainers under emergency policy.
- Protect `v*` tags so release tags cannot be moved or created from unreviewed commits.
- Require maintainer MFA for repository and registry access.
- Document who can create releases and who can rotate release credentials.

### Phase 2: harden dependencies and CI inputs

Outcome: dependencies and CI extensions are visible, reviewed, and updateable.

- Add Dependabot for Go modules, GitHub Actions, Docker base images, and the Hugo UI if package manifests are active.
- Pin third-party GitHub Actions by immutable SHA where compatible.
- Keep SLSA reusable workflows pinned to `@vX.Y.Z` because `slsa-verifier` expects that form.
- Add a lightweight policy for approving new GitHub Actions and external services.
- Add `govulncheck` for Go dependency and standard-library vulnerability checks.
- Add container image scanning for the final image digest.
- Add secret scanning in CI with a tool such as `gitleaks`.

### Phase 3: make artifact integrity enforceable

Outcome: every published artifact can be verified by digest, signature, and provenance.

- Keep GoReleaser checksums attached to every release.
- Keep SBOM generation for release archives and container images.
- Keep SLSA provenance for archives and GHCR images.
- Add a post-release verification job that downloads the just-published assets and runs `slsa-verifier`.
- Add a post-release `cosign verify` check for the GHCR image.

### Phase 4: add release operations and recovery

Outcome: a bad or compromised release has a practiced response path.

- Write a release runbook covering tag creation, release monitoring, verification, and rollback.
- Write a credential rotation checklist for `GITHUB_TOKEN`, Telegram bot token, and publishing tokens.
- Document how to revoke or supersede a bad release.
- Add release audit checks: tag actor, workflow run URL, image digest, checksums, provenance URL, SBOM URL.
- Ensure logs do not print secrets or full authorization URLs.

## Definition of done

The release process is acceptable when a maintainer can answer yes to all of these:

- Did the release tag come from a reviewed commit on `main`?
- Did CI run on GitHub-hosted runners with minimal job permissions?
- Were tests, vulnerability checks, secret scanning, and image scanning completed?
- Are release archives checksummed, SBOMed, signed or attested, and covered by SLSA provenance?
- Is the GHCR image signed, SBOMed, and covered by SLSA provenance?
- Can a consumer verify the archive and image without trusting local build output?
- Is there an incident path for revoking, replacing, or warning about a bad release?

## Near-term work items

Do these next, in order:

1. Add Dependabot coverage for Go modules, GitHub Actions, and Docker.
2. Add `govulncheck` and secret scanning to CI.
3. Add final container image scanning after Buildx produces a digest.
4. Pin non-SLSA third-party Actions by SHA.
5. Add a release verification job using `slsa-verifier` and `cosign verify`.
6. Document GitHub branch and tag protection settings in a release runbook.

## Non-goals for now

- Replacing GitHub Actions with self-hosted build infrastructure.
- Requiring fully hermetic builds.
- Blocking all releases on zero low-severity findings.
