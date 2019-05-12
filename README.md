# XCTestHTMLReport

Generate Xcode-like HTML report for Unit and UI Tests with XCTestHTMLReport

## Current release

![GitHub release](https://img.shields.io/github/release/BirmacherAkos/bitrise-step-xctest-html-report.svg)
![GitHub Release Date](https://img.shields.io/github/release-date/BirmacherAkos/bitrise-step-xctest-html-report.svg)

## CI status

![Bitrise](https://img.shields.io/bitrise/dbb0739f4a28d789.svg?token=HI6D8qe117T1G_O9_Wn9ZQ)
![Codecov](https://img.shields.io/codecov/c/github/BirmacherAkos/bitrise-step-xctest-html-report.svg?token=eeb445314cb94bbaa8ac01bc45cb3d37)

Public CI on Bitrise.io\
https://app.bitrise.io/app/dbb0739f4a28d789#/builds

# How to 
Add this step **after** the **Xcode Test for iOS** step. The XCTestHTMLReport step will search for the `.xcresult` file in the `$BITRISE_XCRESULT_PATH` by default, because the **Xcode Test for iOS** step will generate it there.\
*You can change the search dir by modifying the `test_result_path` step input.*
    
![example_workflow](https://github.com/BirmacherAkos/bitrise-step-xctest-html-report/blob/readme_img_store/readme_img_store/example_workflow.png)

XCTestHTMLReport step will generate the test report files under the `$BITRISE_DEPLOY_DIR`. If you want to make that file available on Bitrise.io add the **Deploy to Bitrise.io - Apps, Logs, Artifacts** step **after** this step.

![example_report](https://github.com/BirmacherAkos/bitrise-step-xctest-html-report/blob/readme_img_store/readme_img_store/example_report.gif)

# Example workflow in bitrise.yml
```
test-simulator-html-report:
    steps:
    - activate-ssh-key:
        run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
    - git-clone: {}
    - cache-pull: {}
    - xcode-test: {}
    - git::https://github.com/BirmacherAkos/bitrise-step-xctest-html-report.git@feature/bitrise_configs:
        inputs:
        - test_result_path: "$BITRISE_XCRESULT_PATH"
    - deploy-to-bitrise-io: {}
    - cache-push: {}
```

---

## How to use a Bitrise Step

Can be run directly with the [bitrise CLI](https://github.com/bitrise-io/bitrise),
just `git clone` this repository, `cd` into it's folder in your Terminal/Command Line
and call `bitrise run test`.

*Check the `bitrise.yml` file for required inputs which have to be
added to your `.bitrise.secrets.yml` file!*

Step by step:

1. Open up your Terminal / Command Line
2. `git clone` the repository
3. `cd` into the directory of the step (the one you just `git clone`d)
5. Create a `.bitrise.secrets.yml` file in the same directory of `bitrise.yml`
   (the `.bitrise.secrets.yml` is a git ignored file, you can store your secrets in it)
6. Check the `bitrise.yml` file for any secret you should set in `.bitrise.secrets.yml`
  * Best practice is to mark these options with something like `# define these in your .bitrise.secrets.yml`, in the `app:envs` section.
7. Once you have all the required secret parameters in your `.bitrise.secrets.yml` you can just run this step with the [bitrise CLI](https://github.com/bitrise-io/bitrise): `bitrise run test`

An example `.bitrise.secrets.yml` file:

```
envs:
- A_SECRET_PARAM_ONE: the value for secret one
- A_SECRET_PARAM_TWO: the value for secret two
```
