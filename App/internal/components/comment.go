package components
var Comment_componnent = `

  <div class="comment">
    <div class="user">
      <div class="profile-photo">
        <img src="/assets/images/4.jpg" />
      </div>
      <div class="info">
       <h3 id="usercomment-{{.Id}}">nom d'utilisateur </h3>
      </div>
    </div>
    <div id="comt-text-style">
      <p>{{.Comment}}</p>
    </div>
    <div class="actions">
      <div class="interaction-buttons">
        <span class="icon"><i class="uil uil-thumbs-up" id="likedcomment-{{.Id}}">{{.Like}}</i></span>
        <span class="icon"><i class="uil uil-thumbs-down"
            id="dislikedcomment-{{.Id}}">{{.Dislike}}</i></span>
            <input type="hidden" id="reactionComment-{{.Id}}" value="">
      </div>
    </div>
  </div>
  
`