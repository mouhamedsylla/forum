package components

var Profil_component = ` <div class="message">
<div class="profile-photo">
  <img src="/assets/images/b.jpg" />
</div>
<div class="message-body">
  <h5>{{.Username}}</h5>
  <br />
  <p class="text-muted">
	<span> <i class="uil uil-pen"></i></span>Publication
  </p>
  <br />
  <p class="text-muted">
	<span> <i class="uil uil-thumbs-up"></i></span> number of like
  </p>
  <br />
  <p class="text-muted">
	<span><i class="uil uil-thumbs-down"></i></span> number of
	dislike
  </p>
  <br />
  <button class="btn" id="btn-deconnecte" >Se déconnecter</button>
</div>
</div>`
