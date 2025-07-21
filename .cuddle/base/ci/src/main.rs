use cuddle_ci::drone_templater::DroneTemplater;
use cuddle_ci::rust_service::architecture::{Architecture, Os};
use cuddle_ci::rust_service::{extensions::*, RustService};
use cuddle_ci::CuddleCI;

const BIN_NAME: &str = "cuddle-go-lib-plan";

#[tokio::main]
async fn main() -> eyre::Result<()> {
    tracing_subscriber::fmt::init();

    dagger_sdk::connect(|client| async move {
        let service = &RustService::from(client.clone())
            .with_arch(Architecture::Amd64)
            .with_os(Os::Linux)
            .with_apt(&[
                "clang",
                "libssl-dev",
                "libz-dev",
                "libgit2-dev",
                "git",
                "openssh-client",
            ])
            .with_apt_release(&["git", "openssh-client"])
            .with_docker_cli()
            .with_cuddle_cli()
            .with_kubectl()
            .with_apt_ca_certificates()
            .with_crates(["ci", "crates/*"])
            .with_mold("2.3.3")
            .with_bin_name(BIN_NAME)
            .with_deployment(false)
            .to_owned();

        let drone_templater = &DroneTemplater::new(client, "templates/cuddle-go-lib-plan.yaml")
            .with_variable("bin_name", BIN_NAME)
            .to_owned();

        CuddleCI::default()
            .with_pull_request(service)
            .with_main(service)
            .with_main(drone_templater)
            .execute(std::env::args())
            .await?;
        Ok(())
    })
    .await?;

    Ok(())
}
