# Making a new release

The release process is now fully automated. Simply run:

```bash
make prepare-release VERSION=v2.37.0
```

This will:
1. Generate and open the changelog for editing
2. Prompt for confirmation to proceed
3. Update version.txt with the new version
4. Create a release branch and PR
5. Wait for the PR to be merged
6. Create and push the git tag
7. Trigger CI to build and publish the release automatically

## Testing releases

To test the release process on your fork:

```bash
make prepare-release VERSION=v2.37.0-test REPO=Mic92/direnv
```
