const mongoHost = process.env.PATIENT_VISIT_API_MONGODB_HOST
const mongoPort = process.env.PATIENT_VISIT_API_MONGODB_PORT

const mongoUser = process.env.PATIENT_VISIT_API_MONGODB_USERNAME
const mongoPassword = process.env.PATIENT_VISIT_API_MONGODB_PASSWORD

const database = process.env.PATIENT_VISIT_API_MONGODB_DATABASE

const COLLECTIONS = [
    "users",
    "patient_visits",
]

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

// try to connect to mongoDB until it is not available
let connection;
while(true) {
    try {
        connection = Mongo(`mongodb://${mongoUser}:${mongoPassword}@${mongoHost}:${mongoPort}`);
        break;
    } catch (exception) {
        console.log(`Cannot connect to mongoDB: ${exception}`);
        console.log(`Will retry after ${retrySeconds} seconds`)
        sleep(retrySeconds * 1000);
    }
}

// if database and collection exists, exit with success - already initialized
const databases = connection.getDBNames()
if (databases.includes(database)) {
    const dbInstance = connection.getDB(database)
    collections = dbInstance.getCollectionNames()
    const existingCollections = [];
    for (const col of COLLECTIONS) {
        if (collections.includes(col)) {
            existingCollections.push(col);
        }
    }

    if (COLLECTIONS.length === existingCollections.length) {
        console.log(`Collections '${COLLECTIONS}' already exist in database '${database}'`)
        process.exit(0);
    }

    if (existingCollections.length > 0) {
        console.log(`Broken state in database '${database}', some collections exist, but not all. Existing collections: ${existingCollections}`)
        process.exit(1);
    }
}

// initialize
// create database and collection
const db = connection.getDB(database)

for (const col of COLLECTIONS) {
    db.createCollection(col)
    // create indexes
    db[col].createIndex({ "id": 1 })
}


//insert sample users data
let result = db["users"].insertMany([
    { "id": "d1", "name": "MUDr. Šefko Doktor", "role": "doctor", "email": "sefko@nemocnica.sk" },
    { "id": "d2", "name": "MUDr. Vedúca Doktorka", "role": "doctor", "email": "veduxa@nemocnica.sk" },
    { "id": "p1", "name": "Lazar Pacientový", "role": "patient", "email": "lazar@email.sk" },
    { "id": "p2", "name": "Dedo Vsevedo", "role": "patient", "email": "dedo@email.sk" },
    { "id": "p3", "name": "Niekto Další", "role": "patient", "email": "niekto@email.sk" },
    { "id": "a1", "name": "Adminus Adminový", "role": "admin", "email": "admin@nemocnica.sk" },
    { "id": "u1", "name": "Ján Novák", "role": "patient", "email": "jan.novak@example.com" },
    { "id": "u2", "name": "Mária Kováčová", "role": "patient", "email": "maria.kovacova@example.com" },
    { "id": "u3", "name": "Peter Horváth", "role": "patient", "email": "peter.horvath@example.com" },
    { "id": "u4", "name": "Lucia Bieliková", "role": "patient", "email": "lucia.bielikova@example.com" },
    { "id": "u5", "name": "Martin Šimek", "role": "patient", "email": "martin.simek@example.com" },
    { "id": "u6", "name": "Anna Králová", "role": "patient", "email": "anna.kralova@example.com" },
    { "id": "u7", "name": "Tomáš Varga", "role": "patient", "email": "tomas.varga@example.com" },
    { "id": "u8", "name": "Zuzana Poláková", "role": "patient", "email": "zuzana.polakova@example.com" },
    { "id": "u9", "name": "Michal Tóth", "role": "patient", "email": "michal.toth@example.com" },
    { "id": "u10", "name": "Eva Marková", "role": "patient", "email": "eva.markova@example.com" },
    { "id": "u11", "name": "Filip Urban", "role": "patient", "email": "filip.urban@example.com" },
    { "id": "u12", "name": "Barbora Sedláčková", "role": "patient", "email": "barbora.sedlackova@example.com" }
]);

if (result.writeError) {
    console.error(result)
    print(`Error when writing the data: ${result.errmsg}`)
}

// exit with success
process.exit(0);
