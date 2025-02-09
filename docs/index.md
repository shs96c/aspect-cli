Aspect is a build system based on Bazel.

The `aspect` CLI is a drop-in replacement for the `bazel` CLI that comes with Bazel.
We also provide Bazel "rules" to integrate tightly with language ecosystems.

Aspect provides a simpler, intutive experience for developers new to Bazel,
adds the power tools that make advanced users more productive
and the customizability your DevInfra team needs to support your organization.

## Customize Bazel for your organization

Every organization has a different engineering culture and developer stack.
Bazel was designed for Google's workflows, not yours.
Aspect's plugin system allows you to fit Bazel into your team's development process,
with custom commands, behaviors, and integrations.

![People working together on software](/people.png)

Your Developer Infrastructure team can write your own plugins using our SDK which execute directly from your repository.
A vibrant ecosystem of plugins will grow which accelerates your Bazel migration.

In the following example, the error message from a badly written `genrule` was confusing,
so a plugin was written to provide more help:

<script id="asciicast-57gaElVKNlb0d8pyZ7JGBDZhL" src="https://asciinema.org/a/57gaElVKNlb0d8pyZ7JGBDZhL.js" async></script>

Some other uses of plugins include:
- stamp out **new Bazel projects** following your local conventions
- point error messages to your **internal documentation**
- **add commands** for deploying, linting, rebasing, or other common developer workflows
- understand where your developers get stuck and **provide help**

Read more in the [plugins documentation](/help/topics/plugins)

## Interactive

When running in an interactive terminal, `aspect` gives helpful prompts to
fix mistakes in how you run the tool, your Bazel configuration, or your code.

In this example, the Bazel configuration didn't allow a dependency because the
`visibility` attribute needed adjustment, so the user is prompted if they'd like
the source files edited:

<script id="asciicast-eL4HL3BZhobRD8U4UIRKzyb8R" src="https://asciinema.org/a/eL4HL3BZhobRD8U4UIRKzyb8R.js" async></script>

Some other examples of interactivity:
- offer to **apply fixes** to your source code that a compiler suggested
- suggest **better usage**, like using `bazel cquery` in place of `bazel query` or avoiding `bazel clean`
- list **common subcommands** or expressions

When run outside an interactive terminal, such as on CI, the prompts are instead printed
for developers to copy-paste to their machine.

## Bazel Rules

We are experts in Bazel rules. In addition to being the maintainers of rules_nodejs and rules_python under the bazelbuild GitHub org, Aspect hosts our own next-generation rulesets.
See the pinned repositories on our GitHub homepage: <https://github.com/aspect-build>

## Expert help is a click away

Our Bazel documentation site, <https://docs.aspect.build> makes it easy to search for
the missing documentation across all Bazel rulesets.

`aspect` is sponsored by Aspect Development, a Bazel consulting company.
If your organization needs more help to make your Bazel migration a success,
come find us at [aspect.dev](https://aspect.dev)

## Open source and no lock-in

You can rely on `aspect` to power your developer experience workflows.

It is a superset of what Bazel provides, so you can always go back to running `bazel` commands.

In fact, it includes `bazelisk` which is the recommended version manager from the Bazel team.

It is open-source, and free for use by individuals, non-profits, and small businesses.
Companies with more than 250 employees can contact us at sales@aspect.dev for a license quote.

# Installation

## Using a package manager

> Coming soon

## Manual installation

Download a binary from our [GitHub Releases] page and put it in your PATH.

On MacOS you can bypass the "Unknown Developer" dialog by running

```shell
xattr -c $(which aspect)
```

before launching `aspect`.

# User Manuals

The commands are documented under [aspect](/aspect).

[Bazel]: http://bazel.build
[GitHub Releases]: https://github.com/aspect-dev/aspect-cli/releases
