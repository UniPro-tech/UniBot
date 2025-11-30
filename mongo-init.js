init = false;
print("Init script ...");

try {
  if (!db.isMaster().ismaster) {
    print("Error: primary not ready, initialize ...");
    rs.initiate({
      _id: "my-replica-set",
      // ここを適切なホスト名に変更
      members: [{ _id: 0, host: "localhost:27017" }],
    });
    quit(1);
  } else {
    if (!init) {
      admin = db.getSiblingDB("admin");
      admin.createUser({
        user: "root",
        pwd: "root",
        roles: ["readWriteAnyDatabase"],
      });
      init = true;
    }
  }
} catch (e) {
  rs.status().ok;
}
