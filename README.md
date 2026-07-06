# GitHub Migration Validator

A GitHub CLI extension for validating GitHub organization and repository migrations by comparing key metrics between source and target repositories. Supports GitHub-to-GitHub, Bitbucket Server (Data Center) to GitHub, and Azure DevOps (Server or Services/cloud) to GitHub migrations.

## Overview

The GitHub Migration Validator helps ensure that your migration to GitHub has been completed successfully. It compares various repository metrics (issues, pull requests, tags, releases, commits) between source and target repositories and provides a detailed validation report.

Supported sources:

- **GitHub** (organization/repository) — full validation including migration archives
- **Bitbucket Server / Data Center** — API-based validation
- **Azure DevOps** (Server on-prem or Services/cloud) — API-based validation (single repo or whole team project)

## Documentation

- **[Migration Archive Support](docs/migration-archive.md)** - Comprehensive guide for enhanced validation using GitHub migration archives

## Install

```bash
gh extension install mona-actions/gh-migration-validator
```

## Usage

### Basic Usage

```bash
gh migration-validator \
  --github-source-org "source-org" \
  --github-target-org "target-org" \
  --source-repo "my-repo" \
  --target-repo "my-repo" \
  --github-source-pat "ghp_xxx" \
  --github-target-pat "ghp_yyy"
```

### With Markdown Output

```bash
gh migration-validator \
  --github-source-org "source-org" \
  --github-target-org "target-org" \
  --source-repo "my-repo" \
  --target-repo "my-repo" \
  --github-source-pat "ghp_xxx" \
  --github-target-pat "ghp_yyy" \
  --markdown-table \
  --markdown-file "validation-report.md"
```

### Skipping LFS Validation

If you want to skip LFS object validation (useful for large repositories or when LFS is not used), use the `--no-lfs` flag:

```bash
gh migration-validator \
  --github-source-org "source-org" \
  --github-target-org "target-org" \
  --source-repo "my-repo" \
  --target-repo "my-repo" \
  --github-source-pat "ghp_xxx" \
  --github-target-pat "ghp_yyy" \
  --no-lfs
```

### Environment Variables

You can use environment variables instead of flags. All environment variables use the `GHMV_` prefix.

#### Source (GitHub-to-GitHub)

```bash
export GHMV_SOURCE_ORGANIZATION="source-org"
export GHMV_SOURCE_TOKEN="ghp_xxx"
export GHMV_SOURCE_REPO="my-repo"
export GHMV_SOURCE_HOSTNAME="https://github.example.com"  # Optional: GitHub Enterprise Server
```

#### Target (shared across all subcommands)

```bash
export GHMV_TARGET_ORGANIZATION="target-org"
export GHMV_TARGET_TOKEN="ghp_yyy"
export GHMV_TARGET_REPO="my-repo"
export GHMV_TARGET_HOSTNAME="https://github.example.com"  # Optional: GitHub Enterprise Server
```

#### Output & Behavior

```bash
export GHMV_MARKDOWN_TABLE="true"                  # Output as markdown table
export GHMV_MARKDOWN_FILE="validation-report.md"   # Write markdown to file
export GHMV_NO_LFS="true"                          # Skip LFS validation
export GHMV_STRICT_EXIT="true"                     # Exit code 2 on validation failures
export GHMV_RATE_LIMIT_THRESHOLD="100"             # GitHub API rate limit warning threshold (default: 50, 0 to disable)
```

#### Basic Run

```bash
gh migration-validator
```

### Strict Exit Mode

Use strict exit mode when you need shell pipelines to detect validation failures. Enable it with the `--strict-exit` flag or set `GHMV_STRICT_EXIT=true` to return exit code 2 whenever any validation fails. Without this option the command exits 0 while still reporting failures in the output.

```bash
gh migration-validator \
  --github-source-org "source-org" \
  --github-target-org "target-org" \
  --source-repo "my-repo" \
  --target-repo "my-repo" \
  --github-source-pat "ghp_xxx" \
  --github-target-pat "ghp_yyy" \
  --strict-exit
```

### GitHub App Authentication

For GitHub App authentication, use environment variables:

```bash
# Source GitHub App
export GHMV_SOURCE_APP_ID="123456"
export GHMV_SOURCE_PRIVATE_KEY="-----BEGIN RSA PRIVATE KEY-----\n..."
export GHMV_SOURCE_INSTALLATION_ID="987654"

# Target GitHub App
export GHMV_TARGET_APP_ID="123457"
export GHMV_TARGET_PRIVATE_KEY="-----BEGIN RSA PRIVATE KEY-----\n..."
export GHMV_TARGET_INSTALLATION_ID="987655"
```

### Enterprise Server Support

For GitHub Enterprise Server:

```bash
export GHMV_SOURCE_HOSTNAME="https://github.example.com"
```

## Export and Validation Workflow

The tool provides both export and validation capabilities that work together to enable point-in-time migration validation:

1. **Export**: Capture repository data at a specific point in time
2. **Validate-from-Export**: Validate target repositories against exported snapshots

This workflow is particularly useful when:

- The source repository continues to receive changes during migration
- You need to validate against the exact state when migration occurred
- You want to create audit trails of migration validation

### Export Usage

```bash
gh migration-validator export \
  --github-source-org "source-org" \
  --source-repo "my-repo" \
  --github-source-pat "ghp_xxx" \
  --format json \
  --output ".exports/my-export.json"
```

### Export with Migration Archive

The tool can also download and analyze migration archives to include additional validation metrics. See the [Migration Archive Documentation](docs/migration-archive.md) for detailed information.

### Export Options

- `--github-source-org` (required): Source organization name
- `--source-repo` (required): Source repository name
- `--github-source-pat` (required): GitHub token with read permissions
- `--source-hostname` (optional): GitHub Enterprise Server URL
- `--format` (optional): Export format - `json` or `csv` (default: `json`)
- `--output` (optional): Output file path (auto-generated if not specified)
- `--download` (optional): Download and analyze migration archive automatically
- `--download-path` (optional): Directory to download migration archives to (default: ./migration-archives)
- `--archive-path` (optional): Path to an existing extracted migration archive directory
- `--no-lfs` (optional): Skip LFS object validation

**Note**: `--download` and `--archive-path` are mutually exclusive. For detailed migration archive usage, see [Migration Archive Documentation](docs/migration-archive.md).

### Export Output Formats

**JSON Format:**

```json
{
  "export_timestamp": "2025-10-13T14:49:08Z",
  "repository_data": {
    "owner": "source-org",
    "name": "my-repo",
    "issues": 42,
    "pull_requests": {
      "open": 5,
      "closed": 10,
      "merged": 15,
      "total": 30
    },
    "tags": 8,
    "releases": 3,
    "commits": 150,
    "latest_commit_sha": "abc123def456",
    "branch_protection_rules": 4,
    "webhooks": 2
  },
  "migration_archive": {
    "issues": 42,
    "pull_requests": 30,
    "protected_branches": 1,
    "releases": 3
  }
}
```

When migration archive data is included, the export will contain additional `migration_archive` metrics. See [Migration Archive Documentation](docs/migration-archive.md) for details.

**CSV Format:**

Contains the same data in CSV format with headers for easy analysis in spreadsheet applications.

### Default Export Location

When no output file is specified, exports are automatically saved to `.exports/` directory with timestamped filenames:

- `.exports/{owner}_{repo}_export_{timestamp}.{format}`

Example: `.exports/mona-actions_my-repo_export_20251002_144908.json`

## Validate-from-Export

The `validate-from-export` command allows you to validate a target repository against a previously exported snapshot of source repository data. This is essential for validating migrations when the source repository may have changed since the migration occurred.

### Validate-from-Export Usage

```bash
gh migration-validator validate-from-export \
  --export-file ".exports/mona-actions_my-repo_export_20251002_144908.json" \
  --github-target-org "target-org" \
  --target-repo "my-repo" \
  --github-target-pat "ghp_yyy"
```

### Using Existing Archive Directory

If you already have an extracted migration archive directory:

```bash
gh migration-validator export \
  --github-source-org "source-org" \
  --source-repo "my-repo" \
  --github-source-pat "ghp_xxx" \
  --archive-path "path/to/extracted/migration-archive"
```

### Validate-from-Export Options

- `--export-file` (required): Path to the exported JSON file containing source data
- `--github-target-org` (required): Target organization name
- `--target-repo` (required): Target repository name
- `--github-target-pat` (required): GitHub token with read permissions for target
- `--target-hostname` (optional): GitHub Enterprise Server URL for target
- `--markdown-table` (optional): Output results in markdown format
- `--markdown-file` (optional): Write markdown output to the specified file; uses the same content without the surrounding ```markdown fences
- `--no-lfs` (optional): Skip LFS object validation
- `--strict-exit` (optional): Exit with status 2 on validation failures

### Environment Variables for Validate-from-Export

```bash
export GHMV_TARGET_ORGANIZATION="target-org"
export GHMV_TARGET_TOKEN="ghp_yyy"
export GHMV_TARGET_REPO="my-repo"
export GHMV_MARKDOWN_TABLE="true"
export GHMV_MARKDOWN_FILE="validation-report.md"
export GHMV_NO_LFS="true"  # Optional: skip LFS validation

gh migration-validator validate-from-export --export-file "path/to/export.json"
```

### Complete Export and Validation Workflow

1. **Export source data before migration:**

   ```bash
   gh migration-validator export \
     --github-source-org "source-org" \
     --source-repo "my-repo" \
     --github-source-pat "ghp_xxx"
   ```

2. **Perform your migration** (using GitHub's migration tools)

3. **Validate against the export:**

   ```bash
   gh migration-validator validate-from-export \
     --export-file ".exports/source-org_my-repo_export_20251002_144908.json" \
     --github-target-org "target-org" \
     --target-repo "my-repo" \
     --github-target-pat "ghp_yyy"
   ```

This ensures you're validating against the exact state of the source repository when the migration occurred, regardless of any subsequent changes.

## Validate Organization

The `validate-org` command validates repositories migrated from a source organization to a target organization in a single run, producing one consolidated report.

By default it lists all repositories in the source org and validates each one against the same-named repo in the target org. This works well for **org-by-org migrations** (e.g. using GEI) where repository names are preserved and all repos (including forks and archives) are migrated. Forked repos have issue validation skipped automatically since forks have issues disabled.

When source and target repository names differ, or you only want to validate a subset, provide a CSV mapping file via `--repo-list`.

### Validate-Org Usage

```bash
gh migration-validator validate-org \
  --github-source-org "source-org" \
  --github-target-org "target-org" \
  --github-source-pat "ghp_xxx" \
  --github-target-pat "ghp_yyy"
```

### With Repository Mapping (CSV)

When not all source repos exist in the target, or names differ:

```csv
# repos.csv
source_repo,target_repo
my-app,my-app-migrated
shared-lib,shared-lib
internal-tools,internal-tools
```

```bash
gh migration-validator validate-org \
  --github-source-org "source-org" \
  --github-target-org "target-org" \
  --github-source-pat "ghp_xxx" \
  --github-target-pat "ghp_yyy" \
  --repo-list repos.csv
```

If a CSV line has only one column, the target name is assumed to match the source. Lines starting with `#` are comments. A header row of `source_repo,target_repo` or `source,target` is automatically skipped.

### With Markdown Report

```bash
gh migration-validator validate-org \
  --github-source-org "source-org" \
  --github-target-org "target-org" \
  --github-source-pat "ghp_xxx" \
  --github-target-pat "ghp_yyy" \
  --markdown-table \
  --markdown-file "org-validation-report.md"
```

When `--markdown-file` is specified, results are written incrementally after each repository completes, so partial progress survives interruptions (e.g. Ctrl+C).

### Environment Variables for Validate-Org

```bash
export GHMV_SOURCE_ORGANIZATION="source-org"
export GHMV_SOURCE_TOKEN="ghp_xxx"
export GHMV_TARGET_ORGANIZATION="target-org"
export GHMV_TARGET_TOKEN="ghp_yyy"

gh migration-validator validate-org
```

### Validate-Org Options

#### Source Flags

- `--github-source-org` / `-s` (required): Source GitHub organization
- `--github-source-pat` / `-a` (required): Source GitHub token with read permissions
- `--source-hostname` / `-u` (optional): GitHub Enterprise Server URL for source

#### Repository Selection

- `--repo-list` (optional): Path to CSV file with `source_repo,target_repo` mappings. When omitted, all repos from the source org are validated with matching names.

#### Shared Target Flags (inherited from root)

- `--github-target-org` / `-t` (required): Target GitHub organization
- `--github-target-pat` / `-b` (required): Target GitHub token with read permissions
- `--target-hostname` / `-v` (optional): GitHub Enterprise Server URL for target
- `--markdown-table` / `-m`: Output results in markdown format
- `--markdown-file`: Write consolidated markdown report to a file
- `--no-lfs`: Skip LFS object validation
- `--strict-exit`: Exit with status 2 on validation failures

### How It Works

1. Determines repositories to validate, either from `--repo-list` CSV or by listing all repos in the source org
2. For each repo pair, validates source against target (names can differ when using CSV mapping)
3. Repos that fail or cannot be accessed are recorded with an error but do not stop validation of remaining repos
4. When `--markdown-file` is set, the report is updated after each repo (incremental writes)
5. Produces a **single summary table** with per-repository pass/fail/warn/info status

### Output

The consolidated report includes:

- A summary table with each repository's overall status (pass/fail/warn/info)
- Total counts across all repositories
- Per-repository detailed validation results (in the markdown report)

## Bitbucket Validation

The `bitbucket` subcommand validates migrations from Bitbucket (Server / Data Center) to GitHub by comparing API metrics between the source Bitbucket instance and the target GitHub repository. This is useful for verifying that repository data was migrated correctly when moving from Bitbucket to GitHub.

### Bitbucket Usage

```bash
gh migration-validator bitbucket \
  --bbs-server-url "https://bitbucket.example.com" \
  --bbs-project "PROJ" \
  --bbs-repo "my-repo" \
  --bbs-token "your-bbs-token" \
  --github-target-org "target-org" \
  --target-repo "my-repo" \
  --github-target-pat "ghp_yyy"
```

### Environment Variables for Bitbucket

```bash
export GHMV_BBS_SERVER_URL="https://bitbucket.example.com"
export GHMV_BBS_PROJECT="PROJ"
export GHMV_BBS_REPO="my-repo"
export GHMV_BBS_TOKEN="your-bbs-token"
export GHMV_TARGET_ORGANIZATION="target-org"
export GHMV_TARGET_TOKEN="ghp_yyy"
export GHMV_TARGET_REPO="my-repo"
gh migration-validator bitbucket
```

### Bitbucket Options

#### Bitbucket Source Flags

- `--bbs-server-url` / `-H` (required): Bitbucket Server URL (aligned with [GEI `bbs2gh`](https://docs.github.com/en/migrations/using-github-enterprise-importer/migrating-from-bitbucket-server-to-github-enterprise-cloud/migrating-repositories-from-bitbucket-server-to-github-enterprise-cloud))
- `--bbs-project` / `-p` (required): Project key (use `~username` for personal repos)
- `--bbs-repo` / `-r` (required): Repository slug
- `--bbs-token` / `-k` (required): Personal access token (or `GHMV_BBS_TOKEN`)

#### Shared Target Flags (inherited from root)

- `--github-target-org` / `-t` (required): Target GitHub organization
- `--github-target-pat` / `-b` (required): Target GitHub token (or `GHMV_TARGET_TOKEN`)
- `--target-repo` (required): Target repository name
- `--target-hostname` / `-v`: GitHub Enterprise Server URL (optional)
- `--markdown-table` / `-m`: Output results in markdown format
- `--markdown-file`: Write markdown output to the specified file
- `--no-lfs`: Skip LFS object validation
- `--strict-exit`: Exit with status 2 on validation failures

### What Gets Validated (Bitbucket → GitHub)

| Metric                                               | Status      | Notes                                         |
| ---------------------------------------------------- | ----------- | --------------------------------------------- |
| Pull Requests (Total, Open, Merged, Declined→Closed) | ✅ Compared | Bitbucket "Declined" maps to GitHub "Closed"  |
| Tags                                                 | ✅ Compared |                                               |
| Commits                                              | ✅ Compared | Default branch only                           |
| Latest Commit SHA                                    | ✅ Compared |                                               |
| Branch Permissions vs Branch Protection Rules        | ℹ️ Advisory | Different concepts — shown for reference only |
| Webhooks                                             | ✅ Compared |                                               |
| Issues                                               | ⏭️ Skipped  | Bitbucket uses Jira, not native issues        |
| Releases                                             | ⏭️ Skipped  | Bitbucket has no equivalent                   |
| LFS Objects                                          | ⏭️ Skipped  | TODO                                          |

### Bitbucket Notes

- Requires Bitbucket Server 5.5+ (uses Bearer token PAT authentication)
- The branch permissions comparison is advisory only (ℹ️ INFO) since Bitbucket branch permissions and GitHub branch protection rules are fundamentally different concepts

## Azure DevOps Validation

The `ado` subcommand validates migrations from Azure DevOps to GitHub by comparing API metrics between the source ADO repository and the target GitHub repository. It works with both **Azure DevOps Server** (on-prem) and **Azure DevOps Services** (cloud). Flag names are aligned with GEI's `ado2gh`.

For cloud (Azure DevOps Services), omit `--ado-server-url` — it defaults to `https://dev.azure.com`. For on-prem (Azure DevOps Server), pass `--ado-server-url` including any collection path prefix (e.g. `https://tfs.example.com/tfs`).

You can validate a **single repository** or an **entire team project** (all repositories) in one run.

### ADO Usage (cloud, single repository)

```bash
gh migration-validator ado \
  --ado-org "my-org" \
  --ado-team-project "MyProject" \
  --ado-repo "my-repo" \
  --ado-pat "your-ado-token" \
  --github-target-org "target-org" \
  --target-repo "my-repo" \
  --github-target-pat "ghp_yyy"
```

### ADO Usage (on-prem, single repository)

```bash
gh migration-validator ado \
  --ado-server-url "https://ado.example.com/tfs" \
  --ado-org "my-collection" \
  --ado-team-project "MyProject" \
  --ado-repo "my-repo" \
  --ado-pat "your-ado-token" \
  --github-target-org "target-org" \
  --target-repo "my-repo" \
  --github-target-pat "ghp_yyy"
```

### ADO Usage (entire team project)

Omit `--ado-repo` to validate every repository in the team project against the same-named repository in the target GitHub organization. This produces a single consolidated report (same format as `validate-org`).

```bash
gh migration-validator ado \
  --ado-server-url "https://ado.example.com/tfs" \
  --ado-org "my-collection" \
  --ado-team-project "MyProject" \
  --ado-pat "your-ado-token" \
  --github-target-org "target-org" \
  --github-target-pat "ghp_yyy" \
  --markdown-file "ado-validation-report.md"
```

When `--markdown-file` is specified in project mode, results are written incrementally after each repository completes, so partial progress survives interruptions.

### ADO Usage (repository name mapping via CSV)

In project mode, when the ADO repository name differs from the target GitHub repository name, or you only want to validate a subset of repositories, provide a CSV mapping file via `--repo-list`. This mirrors the `validate-org` behavior.

```csv
source_repo,target_repo
MyRepo,my-repo-migrated
SharedLib,SharedLib
```

```bash
gh migration-validator ado \
  --ado-org "my-org" \
  --ado-team-project "MyProject" \
  --ado-pat "your-ado-token" \
  --github-target-org "target-org" \
  --github-target-pat "ghp_yyy" \
  --repo-list repos.csv
```

If a line has only one column, the target name is assumed to match the source. Lines starting with `#` are treated as comments. When `--repo-list` is omitted, every repository in the team project is validated against the same-named target repository.

### Environment Variables for ADO

```bash
export GHMV_ADO_SERVER_URL="https://ado.example.com/tfs"  # Optional: omit for Azure DevOps Services (cloud)
export GHMV_ADO_ORG="my-collection"
export GHMV_ADO_TEAM_PROJECT="MyProject"
export GHMV_ADO_REPO="my-repo"        # Optional: omit for whole-project validation
export GHMV_ADO_PAT="your-ado-token"
export GHMV_ADO_API_VERSION="6.0"     # Optional: auto-detected when omitted
export GHMV_TARGET_ORGANIZATION="target-org"
export GHMV_TARGET_TOKEN="ghp_yyy"
export GHMV_TARGET_REPO="my-repo"     # Required only for single-repo validation
gh migration-validator ado
```

### ADO Options

#### ADO Source Flags

- `--ado-server-url` / `-H` (optional): Azure DevOps Server URL for on-prem (e.g., `https://ado.example.com/tfs`). Omit for Azure DevOps Services (cloud); defaults to `https://dev.azure.com`.
- `--ado-org` / `-o` (required): Azure DevOps organization (cloud) or collection (on-prem) name
- `--ado-team-project` / `-P` (required): Azure DevOps team project name
- `--ado-repo` / `-r` (optional): Repository name. Omit to validate all repos in the team project.
- `--ado-pat` / `-k` (required): Azure DevOps personal access token
- `--ado-api-version` (optional): REST API version. Auto-detected when omitted (tries 7.1 down to 4.1). Pin explicitly if needed, e.g. `7.1` for cloud, `4.1`/`5.0` for older TFS/ADO Server.
- `--repo-list` (optional, project mode only): Path to a CSV file with `source_repo,target_repo` mappings. When omitted, all repos in the team project are validated with matching names.

#### Shared Target Flags (inherited from root)

- `--github-target-org` / `-t` (required): Target GitHub organization
- `--github-target-pat` / `-b` (required): Target GitHub token (or `GHMV_TARGET_TOKEN`)
- `--target-repo` (required for single-repo mode): Target repository name
- `--target-hostname` / `-v`: GitHub Enterprise Server URL (optional)
- `--markdown-table` / `-m`: Output results in markdown format
- `--markdown-file`: Write markdown output to the specified file
- `--strict-exit`: Exit with status 2 on validation failures

### What Gets Validated (Azure DevOps → GitHub)

| Metric                                             | Status      | Notes                                             |
| -------------------------------------------------- | ----------- | ------------------------------------------------- |
| Pull Requests (Total, Active→Open, Completed→Merged, Abandoned→Closed) | ✅ Compared | ADO PR states mapped to GitHub terminology        |
| Tags                                               | ✅ Compared |                                                   |
| Commits                                            | ✅ Compared | Default branch only                               |
| Latest Commit SHA                                  | ✅ Compared | Mismatch caused by `git lfs migrate` downgraded to ⚠️ WARN when the parent commit matches (see notes) |
| Branch Policies vs Branch Protection Rules         | ℹ️ Advisory | Different concepts — shown for reference only     |
| Service Hooks                                      | ✅ Compared | Collection-level subscriptions scoped to the repo |
| Issues                                             | ⏭️ Skipped  | ADO uses Work Items, not native git issues        |
| Releases                                           | ⏭️ Skipped  | ADO has no equivalent for git repos               |
| LFS Objects                                        | ✅ Compared | Counted from .gitattributes patterns; skip with `--no-lfs` |

### ADO Notes

- Uses PAT authentication (HTTP Basic). When `--ado-api-version` is omitted, the tool auto-detects a supported REST API version by probing the server (tries newest to oldest: 7.1, 7.0, 6.0, 5.1, 5.0, 4.1), so it works across TFS 2018, ADO Server 2019/2020+, and ADO Services
- To pin a version explicitly, set `--ado-api-version` (or `GHMV_ADO_API_VERSION`); a `404` on an existing repository usually means the pinned version is not supported by that server
- For Azure DevOps Services (cloud), omit `--ado-server-url` (defaults to `https://dev.azure.com`). For on-prem, pass the server URL including any collection path prefix (e.g. `https://tfs.example.com/tfs`); legacy `https://{org}.visualstudio.com` URLs are also supported
- If the server is unreachable, version detection fails fast instead of retrying every candidate version
- LFS objects are counted by reading `.gitattributes` for LFS-tracked patterns and inspecting matching files for pointer content (same approach as the GitHub side); use `--no-lfs` to skip. Note that a manual migration which runs `git lfs migrate` converts previously non-LFS files into LFS, so the target may legitimately report more LFS objects than the ADO source
- When the source has no LFS objects but the target does (a signature of `git lfs migrate` rewriting the tip commit), the Latest Commit SHA mismatch is downgraded from ❌ FAIL to ⚠️ WARN if the parent commit (last commit − 1) still matches on both sides, indicating the history is otherwise identical. The Difference column explains the downgrade. If the parent commits also differ, it remains a ❌ FAIL
- The branch policy comparison is advisory only (ℹ️ INFO) since ADO branch policies and GitHub branch protection rules are fundamentally different concepts
- Service hooks are collection-level subscriptions; the count is a best-effort match on subscriptions scoped to the repository

## Migration Archive Support

The tool supports working with GitHub migration archives for enhanced validation capabilities. Migration archives provide three-way validation comparing Source API ↔ Archive ↔ Target API data.

For comprehensive documentation on migration archive features, workflow, and usage examples, see [Migration Archive Documentation](docs/migration-archive.md).

## What Gets Validated

The tool compares the following metrics between source and target repositories:

- **Issues**: Total count (expects +1 in target for migration log issue)
- **Pull Requests**: Total, Open, Merged, and Closed counts
- **Tags**: Total count of Git tags
- **Releases**: Total count of GitHub releases
- **Commits**: Total commit count on default branch
- **Branch Protection Rules**: Total count of branch protection rules configured for the repository
- **Webhooks**: Total count of active repository webhooks
- **LFS Objects**: Total count of Git LFS (Large File Storage) objects referenced in the repository (can be skipped with `--no-lfs` flag)
- **Latest Commit SHA**: Ensures both repositories have the same latest commit in default branch

## Validation Results

- ✅ **PASS**: Metrics match expected values
- ❌ **FAIL**: Target is missing data from source
- ⚠️ **WARN**: Target has more data than source (usually acceptable)
- ℹ️ **INFO**: Advisory comparison only (e.g., Bitbucket branch permissions vs GitHub branch protection)

## Output Formats

### Console Output

The tool provides a formatted table with colored status indicators and a summary.

Example:

```
📊 Migration Validation Report

🔄 Source vs Target Validation

Metric                                 | Status  | Source Value                             | Target Value                             | Difference
Issues (expected +1 for migration log) | ⚠️ WARN  | 2 (expected target: 3)                   | 7                                        | Extra: 4
Pull Requests (Total)                  | ✅ PASS | 29                                       | 29                                       | Perfect match
Pull Requests (Open)                   | ✅ PASS | 0                                        | 0                                        | Perfect match
Pull Requests (Merged)                 | ✅ PASS | 27                                       | 27                                       | Perfect match
Tags                                   | ✅ PASS | 25                                       | 25                                       | Perfect match
Releases                               | ✅ PASS | 25                                       | 25                                       | Perfect match
Commits                                | ✅ PASS | 64                                       | 64                                       | Perfect match
Branch Protection Rules                | ✅ PASS | 1                                        | 1                                        | Perfect match
Webhooks                               | ✅ PASS | 0                                        | 0                                        | Perfect match
LFS Objects                            | ✅ PASS | 15                                       | 15                                       | Perfect match
Latest Commit SHA                      | ✅ PASS | d11552345ad4ffea894b59d9a4145a5119d77dba | d11552345ad4ffea894b59d9a4145a5119d77dba | N/A

📦 Migration Archive vs Source Validation

Metric                           | Status  | Source API Value | Archive Value | Difference
Archive vs Source Issues         | ❌ FAIL | 2                | 6             | Missing: 4
Archive vs Source Pull Requests  | ✅ PASS | 29               | 29            | Perfect match
Archive vs Source Protected Branches | ✅ PASS | 1            | 1             | Perfect match
Archive vs Source Releases       | ✅ PASS | 25               | 25            | Perfect match

🎯 Migration Archive vs Target Validation

Metric                                                    | Status  | Archive Value | Target Value | Difference
Archive vs Target Issues (expected +1 for migration log)  | ✅ PASS | 6 (expected target: 7) | 7     | Perfect match
Archive vs Target Pull Requests                           | ✅ PASS | 29            | 29           | Perfect match
Archive vs Target Protected Branches                      | ✅ PASS | 1             | 1            | Perfect match
Archive vs Target Releases                                | ✅ PASS | 25            | 25           | Perfect match

📊 Passed: 16  📊 Failed: 1  📊 Warnings: 1

❌ Migration validation FAILED - Some data is missing in target
```

### Markdown Output

Use the `--markdown-table` flag to generate copy-paste ready markdown for documentation.

## Dependencies

- [Go](https://golang.org/doc/install) 1.20 or higher
- Key dependencies:
  - [Cobra](https://github.com/spf13/cobra) - CLI framework
  - [Viper](https://github.com/spf13/viper) - Configuration management
  - [go-github](https://github.com/google/go-github) - GitHub REST API client
  - [githubv4](https://github.com/shurcooL/githubv4) - GitHub GraphQL API client
  - [go-githubauth](https://github.com/jferrl/go-githubauth) - GitHub App authentication
  - [go-github-ratelimit](https://github.com/gofri/go-github-ratelimit) - Rate limit handling
  - [pterm](https://github.com/pterm/pterm) - Terminal styling and formatting

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](.github/contributing.md) for guidelines.

## License

[MIT](./LICENSE) © [Mona-Actions](https://github.com/mona-actions)
