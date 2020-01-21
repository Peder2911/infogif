
// ################################################### 
// DEFINITIONS ####################################### 

let img = document.getElementById("mainGif")
let gif = new SuperGif({
      gif:img,
      loop_mode:false
})

let frameSlider = document.getElementById("frameSlider")
let dropArea = document.getElementById("dropArea")

// ################################################### 
// FUNCTIONS ######################################### 

updateGif = function(src){
   img.setAttribute("src", src)
   img.setAttribute("rel:animated_src", src)
   gif.load(controlGif)
}


controlGif = function(){
   frameSlider.max = gif.get_length()
   frameSlider.value = 1
   frameSlider.oninput = function(){
      gif.move_to(this.value-1);
   }
}

handleFiles = function(files){
   uploadFile(files[0]) 
}

uploadFile = function(file){
   

   let url = 'upload/'
   let formData = new FormData()
   formData.append('file', file)

   fetch(url, {
      method: 'POST',
      body: formData
   })
   .then((response) => {
      return response.text()
   }).then(function(gifhash){
      updateGif("/static/gifs/"+gifhash)
   }).catch(() => {console.log("something went wrong!!")})
}

dropHandler = function(event){
   console.log("Something got dropped");
   let data = event.dataTransfer;
   let files = data.files;
   handleFiles(files)
}

highlight = function(event){
   event.target.classList.add("highlight");
}

unhighlight = function(event){
   event.target.classList.remove("highlight");
}


preventDefaults = function(event){
   event.preventDefault();
   event.stopPropagation();
}

// ################################################### 
// LISTENERS ######################################### 

;['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
  dropArea.addEventListener(eventName, preventDefaults, false);
})

;["dragenter","dragover"].forEach(eventName => {
   dropArea.addEventListener(eventName, highlight, false);
})


;["dragleave","drop"].forEach(eventName => {
   dropArea.addEventListener(eventName, unhighlight, false);
})

dropArea.addEventListener("drop", dropHandler, false);

