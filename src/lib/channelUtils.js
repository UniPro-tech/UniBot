const GetLogChannel = async (interaction) => {
  const channel = await interaction.client.channels
    .fetch(interaction.client.conf.logch.command)
    .catch((error) => null);
  return channel;
};

const GetErrorChannel = async (interaction) => {
  const channel = await interaction.client.channels
    .fetch(interaction.client.conf.logch.error)
    .catch((error) => null);
  return channel;
};

module.exports = { GetLogChannel, GetErrorChannel };
