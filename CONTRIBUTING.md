# Making a new release

- `export version=v<version>`
- `git changelog` to generate the CHANGELOG.md, do a bit of editing there.
- Update the version.txt
- `git add version.txt CHANGELOG.md`
- `git commit -m "Release $version"`
- Create a pull request for that and merge it.
- `git tag $version`
- `git push origin $version`
- Wait for the draft release created by ci trigged by pushing the tag
- `make dist`
- `gh release upload $version dist/direnv.*`
- Click the release button
