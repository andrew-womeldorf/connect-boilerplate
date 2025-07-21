#!/usr/bin/env python3
#
# mocking https://github.com/localstack/localstack-extensions/tree/main/terraform-init

import logging
import os

from pip import main as pip

LOG = logging.getLogger(__name__)


def install_localstack_packages():
    pip(["install", "localstack"])
    pip(["install", "terraform-local"])


def install_terraform():
    from localstack.packages.terraform import terraform_package

    terraform_package.install()


def terraform_apply():
    from localstack.packages.terraform import terraform_package
    from localstack.utils.run import run

    tf_path = terraform_package.get_installed_dir()
    # install_dir = tflocal_package.get_installer()._get_install_dir(
    # InstallTarget.VAR_LIBS
    # )
    # tflocal_path = f"{install_dir}/bin"
    env_path = f"{tf_path}:{os.getenv('PATH')}"
    # env_path = f"{tflocal_path}:{tf_path}:{os.getenv('PATH')}"

    workdir = "/terraform"
    LOG.info("Applying terraform project from file %s", workdir)
    # run tflocal
    # workdir = os.path.dirname(path)
    LOG.debug("Initializing terraform provider in %s", workdir)
    run(
        ["tflocal", f"-chdir={workdir}", "init", "-input=false"],
        env_vars={"PATH": env_path},
    )
    LOG.debug("Applying terraform file %s", workdir)
    run(
        ["tflocal", f"-chdir={workdir}", "apply", "-auto-approve"],
        env_vars={"PATH": env_path},
    )


def main():
    install_localstack_packages()
    install_terraform()
    terraform_apply()


if __name__ == "__main__":
    main()
