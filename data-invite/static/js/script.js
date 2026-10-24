const noBtn=document.getElementById('no-btn');
noBtn.addEventListener('mouseover',()=>{
    const x=Math.random()*(window.innerWidth-noBtn.offsetWidth);
    consty=Math.random()*(window.innerHeight-noBtn.offsetHeight);
    noBtn.style.position='absolute';
    noBtn.style.left='${x}px';
    noBtn.style.top='${y}px';
});