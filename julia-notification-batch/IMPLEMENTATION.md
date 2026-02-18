# Documentazione Tecnica - Julia Notification Batch

Questo documento descrive l'architettura e l'implementazione del servizio Batch Go per l'orchestrazione delle preferenze di notifica.

## Architettura
Il progetto è un servizio **Job Scheduler** che esegue periodicamente compiti di sincronizzazione tra API esterne e interne.

### Componenti Principali

#### 1. Scheduler (`internal/scheduler`)
- **Perché**: Permette l'esecuzione temporizzata seguendo la sintassi Cron.
- **Come**: Basato su un motore cron che gestisce il ciclo di vita dei job registrati.

#### 2. Notification Job (`internal/jobs`)
- **Perché**: Incapsula la logica di esecuzione di un'unità di lavoro.
- **Come**: Implementa una logica di **Retry/Refire** simile a quella di Quartz (Java). Se l'orchestrazione fallisce, il job tenta nuovamente l'esecuzione fino a un numero massimo di tentativi configurato (`maxRetries`).

#### 3. Orchestrator Service (`internal/service`)
- **Perché**: Coordina il flusso di dati tra diversi sistemi senza accoppiarli tra loro.
- **Come**: Utilizza delle interfacce (`Gateway`) per interagire con i servizi esterni. Questo permette una facile testabilità tramite mock.
- **Mapping**: Gestisce la trasformazione (mapping) tra i modelli ricevuti dalle API esterne e quelli richiesti dalle API interne.

#### 4. Gateways (`internal/gateway`)
- **Perché**: Astrazione dei client REST.
- **Come**: Implementano chiamate HTTP sincrone con timeout configurabili.
    - `ExternalNewsGateway`: Recupera le notifiche/news da un provider esterno.
    - `NotificationNewsGateway`: Aggiorna le preferenze nel sistema di notifica interno tramite metodo `PUT`.

## Flusso di Orchestrazione
1. **Trigger**: Il `Scheduler` avvia il `NotificationJob` in base alla cron expression (configurata in `application.yaml`).
2. **Recupero**: L'orchestratore interroga `ExternalNewsGateway` per ottenere la lista delle news correnti.
3. **Trasformazione**: I dati vengono normalizzati filtrando i campi necessari.
4. **Sincronizzazione**: L'orchestratore invia i dati normalizzati al `NotificationNewsGateway`.
5. **Monitoraggio**: Ogni fase viene loggata per permettere il tracciamento delle performance e degli errori.
