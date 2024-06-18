document.addEventListener('DOMContentLoaded', function () {
    document.querySelector('button').addEventListener('click', async () => {
        var file = document.getElementById('docfile').files[0]
        if (!file) {
            return
        }

        let [tab] = await chrome.tabs.query({ active: true, currentWindow: true })
        
        var reader = new FileReader()
        reader.onload = function(e) {
          var j = JSON.parse(e.target.result)
          chrome.scripting.executeScript({
            target: { tabId: tab.id },
            func: (j) => {
                window['data'].encounters = j.encounters
                window['data'].encounter_pools = j.encounter_pools
                window['encounterPools'] = j.encounter_pools

                pokemonEncounters.clear()
                for (let i in j.encounters) {
                  addPoolInfo(j.encounters[i])
                }
            },
            args: [j],
            world: 'MAIN',
          })
        };
        reader.readAsText(file)
    });
});