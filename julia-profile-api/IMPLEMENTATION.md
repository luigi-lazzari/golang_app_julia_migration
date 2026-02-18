# Documentazione Tecnica - Julia Profile BFF (api)

Questo documento descrive l'implementazione del Backend-for-Frontend (BFF) dedicato alla gestione del profilo utente e delle preferenze.

## Architettura
Il servizio gestisce dati sensibili e preferenze utente, interfacciandosi con sistemi di persistenza (Cosmos DB), caching (Redis) e servizi di notifica.

### Componenti Principali

#### 1. User Profile Service (`internal/service`)
- **Perché**: Gestisce il recupero e l'aggiornamento dei dati anagrafici dell'utente.
- **Caching (Redis)**: Implementa una strategia di caching **Read-Through**.
    - Se il profilo è in cache (Redis), viene restituito immediatamente.
    - Se non presente, viene recuperato dal database (Cosmos DB) e salvato in cache con un TTL configurabile.
    - All'aggiornamento del profilo, la cache viene invalidata (`Delete`) per garantire la consistenza dei dati.

#### 2. User Preferences Service (`internal/service`)
- **Perché**: Gestisce le preferenze relative alla chat, alla lingua e alle notifiche push.
- **Merge Logic**: Quando l'utente richiede le preferenze, il servizio fonde (merge) i valori di default definiti nel file di configurazione con le scelte effettuate dall'utente e salvate nel database.
- **Sincronizzazione**: Quando un'utente aggiorna le preferenze di chat, il servizio avvia una goroutine (**Fire-and-Forget**) per sincronizzare queste scelte con il `Notification Service` attraverso un client interno.

#### 3. Notification Client (`internal/client`)
- **Perché**: Permette la comunicazione con il microservizio Julia Notification.
- **Come**: Gestisce la registrazione/cancellazione delle installazioni dei dispositivi (`UpsertInstallation`) per permettere l'invio corretto delle notifiche push.

#### 4. Repository layer (`internal/repository`)
- **Perché**: Astrazione dell'accesso ai dati.
- **Come**: Definisce interfacce per Cosmos DB, permettendo al business layer di rimanere agnostico rispetto alla tecnologia di persistenza.

## Logiche di Business
- **Multi-Piattaforma**: Gestisce identificativi differenti per le piattaforme Android e iOS nel sistema di preferenze.
- **Custom Preferences**: Supporta l'aggiunta di descrizioni personalizzate per specifiche preferenze utente (es. preferenze chat estese).

## Stato dell'Integrazione
Il servizio è predisposto per l'uso di **Redis** e **Cosmos DB**. Molte chiamate al database sono attualmente simulate tramite mock in attesa della configurazione definitiva degli endpoint di produzione.
