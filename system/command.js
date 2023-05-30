module.exports = {
    async handling(client, fs, Collection, config) {
        // コマンドハンドリング
        client.commands = new Collection();
        const commandFolders = fs.readdirSync(`./commands`);
        for (const folder of commandFolders) {
            console.log(`\u001b[32m===${folder} commands===\u001b[0m`);
            const commandFiles = fs.readdirSync(`./commands/${folder}`).filter(file => file.endsWith(".js"));
            for (const file of commandFiles) {
                console.log(`dir:${folder},file:${file}`);
                try{
                    const command = require(`../commands/${folder}/${file}`);
                }
                catch{
                    console.log("error!!");
                }
                try {
                    client.commands.set(command.data.name, command);
                    console.log(`${command.data.name} がロードされました。`);
                } catch (error) {
                    console.log(`\u001b[31m${command.data.name} はエラーによりロードされませんでした。\nエラー内容\n ${error}\u001b[0m`);
                }
            }
            console.log(`\u001b[32m===${folder} loaded===\u001b[0m`);
        }
    }
}