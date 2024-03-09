package components

var Post_component = `
<div class="head">
                <div class="user">
                  <div class="profile-photo">
                    <img src="/assets/images/ND.jpg" />
                  </div>
                  <div class="info">
                    <h3 id="name-{{.Id}}"></h3>
                    <small>sénégal, {{.CreatedAt}}</small>
                  </div>
                </div>
              </div>
              <div class="photo">
                {{.Content}}
              </div>
              <div class="action-buttons">
                <div class="interaction-buttons">
                  <span class="icon"><i class="uil uil-thumbs-up" id="liked-{{.Id}}">{{.Like}}</i></span>
                  <span class="icon"><i class="uil uil-thumbs-down" id="disliked-{{.Id}}">{{.Dislike}}</i></span>
                  <input type="hidden" id="reaction-{{.Id}}" value="">
                </div>
                <div class="bookmark">
                  <span><i class="uil uil-comment-dots" id="comment-{{.Id}}"></i></span>
                  <input type="hidden" id="state-{{.Id}}" value="close">
                </div>
              </div>
              <div class="liked-by">
                <span><img src="/assets/images/b.jpg" /></span>
                <span><img src="/assets/images/ND.jpg" /></span>
                <p>Liked by <b>fatou diagne</b> and <b>2, 323 others</b></p>
              </div>
              <div class="caption">
                <p id="typeCategorie-{{.Id}}">
                  
                </p>
              </div>
`
