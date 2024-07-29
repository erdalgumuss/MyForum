# Working On:




# Yapilacaklar:
<br>Guest iken konu olusturmuyo ama hata dondurmuyo, login pop up'ini acsin
Register basarili ise alert versin.<br/><br/>
<br>Ana Sayfa ayarlansin
<br>
<br>

# CSS problems:
~~the structure dissappearing in the footer.~~<br/>
body { min-height } ile ayarlandi<br/><br/>
header'in ayrilmasi gerekiyo, belki iyilestirilebilir.<br/>
base de (main-container) olayinin ayarlanmasi gerekiyo<br/>
suan main-container da ki img headera da yansiyo
<br/>
> Profilde::: <br><br>
> kullanici adi<br>
> email<br>
> Limit ? kac tane gosterilecek ve `show more`<br>
> EKLENSIN
<br>
<br>

**JS:**

**pop-up / main.js icerisinden ayrilabilir**<br/>
resim 20 mbden buyuk oldugunda RAW hatasi donduruyo alert dondurmeli<br/>
Forumda konu olmadiginda error fetching hatasi donduruyo<br/><br/><br><br>

# Tidy:
forum_controller/getPosts ve forum_controller/getPosts -- `silindi`<br/>
modes/topic.go ihtiyac var mi ? -- `yorum satirina alindi`<br/>
likePost, dislikePost, likeComment, dislikeComment -- `controllersdan silindi`<br/>
handlers'dan kullaniliyor<br/>
# Dublicating Stuff:
GetUserProfile in handlers/profile.go and in controller profile_controller.go<br/>
createPost<br/><br/><br><br>
logout dublicating in auth_controller and auth.go<br/>

# QUESTION:
Controllers neden gerekli ? cagirmasi mi gerekli yoksa ikincil func. mi gerekli ?<br/>
script src ???? base de olmasi yeterli mi ? hangi dosyada hangi import olmasi gerekli ?<br/>
{{base}} olayi incelenmeli, bi tik degisti, yorum satirina alinan var profile.html de<br/>
**script src MAIN.js 4 htmlde de vardi. suan sadece base de var. herseyi 2 kez yapma cozuldu, dublicate register cozuldu**<br><br>
CreatePost/GetPost ve CreateComment/GetComment func. time code is different but working fine ???<br/>

# Future: #
<br>

**MANAGEMENT:**<br/>
Admin paneli<br/>
Moderator paneli<br/>
**For Staff:**<br/>
Post delete<br/>
Comment delete<br/>
**For User:**<br/>
Mesaj edit<br/>
Comment edit<br/>
<br>
Report/Bildir<br/>

**FORUM:**<br/>
Mesajlasma (PM/DM)
<br>
<br>

# ACHIEVED / SUCCESSED :: 

<br>~~konu olusturulduktan sonra, olusturulan konuya yonlendirmeli~~
<br>~~comment icinde oyle, **comment bunu suan go uzerinden yapiyo.**~~
<br>~~Post/Comment (created_at) birbirinden farkli, yanlis gosteriyo~~
<br>~~Bir cok islem 2 kez gerceklestiriliyor: **suan kullaniciyi db de 2 kez kaydediyor**~~
<br>~~(like_dislike) seperated from main.js to like_dislike.js~~
<br>~~JS has both on profile.html(injected) AND main.js~~
<br>~~User likes on profile page~~
<br>~~User comments on profile page~~
<br>~~User registered(created_at)~~
<br>~~Filtering~~<br/>
<br>
Foruma suan akla gelmeyen, duzenlenirken elde edilen<br/>
ve bazen 0dan olmayan, gin kullaniyo olunca uygulamasi zor olan<br/>
bu da kalsin dedirten bir cok (ozellikle gercek bir forumda) olan<br/>
nitelikler mevcut. Test, bilme, yapiyi anlama, kafada olusturma<br/>
neyin ne oldugunu, ne kadar profesyonelce oldugunu gorme adina::<br/>
sonradan yapilan fixlerin, **OZELLIKLE ESKI README LERE BAKILARAK**<br/>
buraya eklenmesi gerektigini dusunuyorum<br/>

## MyForum

SQL MÜQ bro

<br>**MyForum/**
<br>├── backend
<br>│   ├── config
<br>│   │   └── config.go
<br>│   ├── controllers
<br>│   │   ├── admin_controller.go
<br>│   │   ├── auth_controller.go
<br>│   │   ├── forum_controller.go
<br>│   │   ├── moderator_controller.go
<br>│   │   └── profile_controller.go
<br>│   ├── DockerFile
<br>│   ├── forum.db
<br>│   ├── go.mod
<br>│   ├── go.sum
<br>│   ├── handlers
<br>│   │   ├── auth.go
<br>│   │   ├── forum.go
<br>│   │   └── profile.go
<br>│   ├── main.go
<br>│   ├── models
<br>│   │   ├── category.go
<br>│   │   ├── comment.go
<br>│   │   ├── like.go
<br>│   │   ├── post.go
<br>│   │   ├── profile.go
<br>│   │   ├── session.go
<br>│   │   ├── topic.go
<br>│   │   └── user.go
<br>│   ├── routes
<br>│   │   ├── admin_routes.go
<br>│   │   ├── auth_routes.go
<br>│   │   ├── forum_routes.go
<br>│   │   ├── moderator_routes.go
<br>│   │   └── profile_routes.go
<br>│   └── utils
<br>│       └── utils.go
<br>├── docker-compose.yml
<br>├── forum.db-x-users-5-password.bin
<br>├── frontend
<br>│   ├── static
<br>│   │   ├── css
<br>│   │   │   └── style.css
<br>│   │   ├── favicon.ico
<br>│   │   ├── images
<br>│   │   │   ├── default-profile.png
<br>│   │   │   ├── hells.jpg
<br>│   │   │   ├── never_stop_riding.jpg
<br>│   │   │   ├── resized.jpg
<br>│   │   │   ├── soa2.jpg
<br>│   │   │   ├── soa3.jpg
<br>│   │   │   ├── soa.jpg
<br>│   │   │   ├── Sonny_Barger.jpg
<br>│   │   │   └── wp.jpg
<br>│   │   └── js
<br>│   │       ├── like_dislike.js
<br>│   │       ├── main.js
<br>│   │       └── profile.js
<br>│   ├── templates
<br>│   │   ├── admin_dashboard.html
<br>│   │   ├── admin.html
<br>│   │   ├── base.html
<br>│   │   ├── comment.html
<br>│   │   ├── edit_post.html
<br>│   │   ├── forum.html
<br>│   │   ├── gallery.html
<br>│   │   ├── index.html
<br>│   │   ├── pending_posts.html
<br>│   │   ├── post.html
<br>│   │   ├── profile.html
<br>│   │   ├── request_moderator.html
<br>│   │   ├── rules.html
<br>│   │   └── user_profile.html
<br>│   └── uploads
<br>│       ├── _117310488_16.jpg
<br>│       ├── 697b023b-64a5-49a0-8059-27b963453fb1.gif
<br>│       ├── 6c0eb42899de8820e8e699d42285e107.jpg
<br>│       ├── 9D798CBA-D927-433B-A11E-FAD76E4C96AF.JPEG
<br>│       ├── media-1576532915.jpeg
<br>│       ├── Screenshot from 2024-07-07 18-07-09.png
<br>│       ├── WIN_20220912_02_15_52_Pro.jpg
<br>│       ├── WIN_20220912_02_15_59_Pro.jpg
<br>│       ├── WIN_20231009_16_47_03_Pro.jpg
<br>│       ├── WIN_20231012_01_56_53_Pro.jpg
<br>│       └── WIN_20231018_16_45_51_Pro.jpg
<br>└── README.md

15 directories, 70 files