use async_trait::async_trait;
use cuddle_ci::{cuddle_please, Context, CuddleCI, MainAction, PullRequestAction};
use dagger_sdk::HostDirectoryOptsBuilder;

#[tokio::main]
async fn main() -> eyre::Result<()> {
    dagger_sdk::connect(|client| async move {
        let service = &GoLib {
            client: client.clone(),
        };
        let cuddle_please = &cuddle_please::CuddlePlease::new(client.clone());

        CuddleCI::default()
            .with_pull_request(service)
            .with_main(service)
            .with_main(cuddle_please)
            .execute(std::env::args())
            .await?;

        Ok(())
    })
    .await?;
    Ok(())
}

#[derive(Clone)]
struct GoLib {
    client: dagger_sdk::Query,
}

impl GoLib {
    pub async fn test(&self) -> eyre::Result<()> {
        let base = self.client.container().from("golang");

        base.with_workdir("/app")
            .with_directory(
                ".",
                self.client.host().directory_opts(
                    ".",
                    HostDirectoryOptsBuilder::default()
                        .include(vec!["**/go.mod", "**/go.sum"])
                        .build()?,
                ),
            )
            .with_exec(vec!["go", "mod", "download"])
            .with_directory(
                ".",
                self.client.host().directory_opts(
                    ".",
                    HostDirectoryOptsBuilder::default()
                        .include(vec!["**/go.mod", "**/go.sum"])
                        .build()?,
                ),
            )
            .with_exec(vec!["go", "test", "./..."])
            .sync()
            .await?;

        Ok(())
    }
}

#[async_trait]
impl PullRequestAction for GoLib {
    async fn execute_pull_request(&self, _ctx: &mut Context) -> eyre::Result<()> {
        self.test().await?;

        Ok(())
    }
}

#[async_trait]
impl MainAction for GoLib {
    async fn execute_main(&self, _ctx: &mut Context) -> eyre::Result<()> {
        self.test().await?;

        Ok(())
    }
}
