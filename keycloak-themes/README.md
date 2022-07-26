# grants-themes

A collection of custom themes for Keycloak

## Dependencies

1. OpenJDK
2. Maven
3. NPM

## How to develop

In general, follow [this guide](https://www.keycloak.org/docs/latest/server_development/#creating-a-theme). You can download Keycloak distribution from the official website,
put theme files into `themes` directory, start a dev server with `./bin/kc.sh start-dev`, and your themes
will be available for use without caching or building anything.

## How to build

1. Run `mvn package`. This will create 2 JAR files in `./target`.
2. Copy the JAR without `-sources` in the name into Keycloak,
putting it into `/opt/jboss/keycloak/standalone/deployments`.
3. If done correctly, you should see the following lines in the logs
```
INFO  [org.keycloak.theme.DefaultThemeManagerFactory] (MSC service thread 1-1) Cleared theme cache
INFO  [org.jboss.as.server] (DeploymentScanner-threads - 1) WFLYSRV0010: Deployed "keycloak-themes-999-SNAPSHOT.jar" (runtime-name : "keycloak-themes-999-SNAPSHOT.jar")
```
4. Your theme should now appear in the list of themes in Realm settings / Themes tab

