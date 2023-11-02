async function waitForLoad() {
    return new Promise((resolve, reject) => {
        if (document.readyState === 'complete') {
            resolve()
        } else {
            window.addEventListener('load', () => {
                resolve()
            })
        }
    })
}

(async () => {
    await waitForLoad()

    const button = document.getElementById('fileupload')
    const chooser = document.querySelector('input[type="file"]');
    const keyTextArea = document.getElementById('key')
    const fileList = document.getElementById('file_list')

    let files = {}

    button.onclick = () => {
        chooser.click()
    }

    keyTextArea.onchange = () => {
        localStorage.setItem('key', keyTextArea.value)
    }

// once localstorage is set, set the value of the key text area
    if (localStorage.getItem('key')) {
        keyTextArea.value = localStorage.getItem('key')
    }

    function addFileToList(fileHash) {
        const file = files[fileHash]
        const li = document.createElement('el')
        li.innerHTML = `
<div class="relative bg-emerald-600 m-5 p-7 rounded-2xl">
<a class="absolute inset-x-1 p-3 top-0 font-medium tracking-tight text-white" href="${window.location.href}cdn/${fileHash}">${file.filename}</a>
<button id="delete-${fileHash}" class="absolute inset-y-0 right-0 p-3 m-1 font-medium tracking-tight text-white bg-red-700 rounded-3xl transition-all ease-in-out hover:bg-red-600">Delete</button>
</div>`
        li.querySelector(`#delete-${fileHash}`).onclick = async () => {
            await fetch(`${window.location.href}delete/${fileHash}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': localStorage.getItem('key')
                }
            })
            await syncFiles()
        }
        fileList.appendChild(li)
    }

    async function syncFiles() {
        const resp = await fetch(`${window.location.href}total_files`, {
            headers: {
                'Authorization': localStorage.getItem('key')
            }
        })
        fileList.innerHTML = ''
        const data = await resp.json()
        files = data.files
        Object.keys(data.files).forEach((file) => {
            addFileToList(file)
        })
    }

    syncFiles()

    chooser.addEventListener('change', async (evt) => {
        const formData = new FormData();
        const file = evt.target.files[0]
        formData.append('file', file, file.name);
        button.innerHTML = `Uploading ${file.name}...`
        // Post request to server
        const resp = await fetch(`${window.location.href}upload`, {
            method: 'POST',
            body: formData,
            headers: {
                'Authorization': localStorage.getItem('key')
            }
        }).catch((err) => {
            button.innerHTML = `File Upload Failed! Server offline or key wrong?`
            console.error(err)
            setTimeout(() => {
                button.innerHTML = `Upload File`
            }, 2500)
        })
        const data = await resp.json()
        await navigator.clipboard.writeText(`${window.location.href}cdn/${data.hash}`)
        syncFiles()
        button.innerHTML = `Copied to clipboard!`
        setTimeout(() => {
            button.innerHTML = `Upload File`
        }, 2500)
    }, false);
})();