# Documentazione Tecnica - Julia Notification Worker

Questo documento descrive l'architettura e l'implementazione del worker Go per la gestione delle notifiche push.

## Architettura
Il progetto segue il pattern **Consumer**, processando messaggi provenienti da Azure Service Bus e inviando notifiche push tramite Azure Notification Hub.

### Componenti Principali

#### 1. Service Bus Client (`internal/servicebus`)
- **Perché**: Gestisce la connessione a bassa latenza con Azure.
- **Come**: Utilizza l'SDK `azservicebus` v1.x. Implementa un loop di ricezione che utilizza la modalità **PeekLock**. Questo garantisce che un messaggio non vada perso se il worker crasha durante l'elaborazione.
- **Settlement**: Il messaggio viene confermato (`Complete`) solo se l'invio all'Hub ha successo, altrimenti viene rilasciato (`Abandon`) per un nuovo tentativo.

#### 2. Service Bus Notification Processor (`internal/worker`)
- **Perché**: Separa la logica di business dall'infrastruttura di trasporto.
- **Come**: Riceve il payload grezzo, lo deserializza in un DTO e lo mappa nel modello di dominio. Gestisce la validazione dei campi obbligatori (`title`, `body`).
- **Deduplicazione**: Intercetta errori di tipo `DuplicateMessageError` per evitare che Service Bus consideri l'elaborazione fallita quando un messaggio è già stato processato.

#### 3. Notification Hub Service (`internal/service`)
- **Perché**: Azure non fornisce un SDK Go aggiornato per Notification Hub, quindi è stata implementata una client REST.
- **Come**: Implementa l'invio di **Template Notifications**. Genera programmaticamente i token **SAS (Shared Access Signature)** necessari per l'autenticazione HMAC-SHA256.
- **Tag**: Supporta `TagExpression` tramite l'header `ServiceBusNotification-Tags` per l'invio targetizzato.

#### 4. Deduplication Service (`internal/service`)
- **Perché**: Garantisce che l'utente non riceva notifiche doppie in caso di retry transienti di Service Bus.
- **Come**: Utilizza una `map` thread-safe (`sync.RWMutex`) in memoria. Ogni messaggio processato con successo viene memorizzato con un **TTL di 24 ore**.
- **Cleanup**: Una goroutine in background esegue la pulizia delle entry scadute ogni 5 minuti per prevenire leak di memoria.

## Flusso di Elaborazione
1. **Ricezione**: Il `servicebus.Client` preleva un messaggio (PeekLock).
2. **Preprocessing**: Il `Processor` deserializza il JSON e applica logiche di fallback (es. `message` -> `body`).
3. **Controllo Duplicati**: Il `DeduplicationService` verifica se il `messageId` è già presente.
4. **Invio**: Se nuovo, il `NotificationHubService` invia la richiesta POST all'Hub con il SAS Token aggiornato.
5. **Conferma**: Se l'invio ha successo, il messaggio viene marcato come processato e confermato su Service Bus.
