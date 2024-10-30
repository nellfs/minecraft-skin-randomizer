# Minecraft Skin Randomizer

This program allows you to mix and match different skin parts to create a new random Minecraft skin.

To run the Minecraft Skin Randomizer, you'll need to provide a configuration file, to specify your configuration file, run the randomizer with the config flag:

```--config=/path/to/your/config/config.json```

 You can find an example config [HERE](https://github.com/nellfs/minecraft-skin-randomizer/blob/main/example/config.json).

---

## Configuration File


```
"edit_skin": "/path/mcpelauncher/games/com.mojang/custom_skins/skin_to_replace.png",
```
I created this software with Minecraft Bedrock in mind (yes, it works for Java skins as well). 
Minecraft Bedrock saves skins locally, so I designed it so that `edit_skin` would be the location where the new skin is generated.
The idea is that you select a custom skin in the game and use it to replaced by the generated random skin.

```
"randomizer_folder": "/path/example/folder",
```

This is the folder used to mix and generate a random skin based on the combination of skins. The more skins you use, the more combinations you'll have. 
You can customize every part of the skin; just place each part in its corresponding folder. 
Randomizer folder template [HERE](https://github.com/nellfs/minecraft-skin-randomizer/tree/main/example/skins).
