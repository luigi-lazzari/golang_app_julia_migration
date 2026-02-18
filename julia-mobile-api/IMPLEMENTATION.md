# Documentazione Tecnica - Julia Mobile BFF (api)

Questo documento descrive l'implementazione del Backend-for-Frontend (BFF) dedicato all'applicazione mobile.

## Architettura
Il servizio funge da aggregatore di configurazioni e stato del sistema per i client mobile (iOS/Android).

### Componenti Principali

#### 1. AppConfig Service (`internal/service`)
- **Perché**: Centralizza la logica di distribuzione delle configurazioni e il controllo delle versioni.
- **Come**: Gestisce la risposta `AppConfigResponse` che include il tempo del server, lo stato di manutenzione e le policy di aggiornamento.
- **Versionamento (SemVer)**: Utilizza la libreria `semver/v3` per confrontare la versione dell'app installata con le versioni minima e raccomandata definite lato server.
    - `REQUIRE`: Se la versione è inferiore alla minima (aggiornamento forzato).
    - `RECOMMEND`: Se la versione è inferiore all'ultima ma superiore alla minima.
    - `NONE`: Se l'app è aggiornata.

#### 2. Configurazione (`internal/config`)
- **Perché**: Permette di definire valori di default per le diverse piattaforme.
- **Come**: Carica impostazioni tramite Viper. Include mappe per `Config`, `Locale` e `Features` che vengono restituite al client per abilitare/disabilitare funzionalità dinamicamente.

#### 3. Handlers (`internal/handler`)
- **Perché**: Espone gli endpoint REST.
- **Come**: Valida gli header obbligatori (`X-App-Platform`, `X-App-Version`) e delega la logica al service layer.

## Logiche di Business
- **Manutenzione**: Il servizio può restituire uno stato di manutenzione (`Enabled: true`) che istruisce l'app a mostrare una schermata di blocco, suggerendo un tempo di retry.
- **Dynamic Features**: Attraverso la sezione `Features` della configurazione, il BFF può abilitare funzionalità (es. nuovi moduli chat o mappe) senza richiedere un rilascio dell'app negli store.

## Stato dell'Integrazione
Attualmente, l'integrazione con **Azure App Configuration** è predisposta ma utilizza ancora dei mock nei service per i valori di versione e URL degli store, in attesa del completamento del gateway repository.
