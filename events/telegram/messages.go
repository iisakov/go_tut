package telegram

const msgHelp = `Этот бот сохраняет ссылки на статьи и умеет выдавать их в рандомном порядке.
Если вам нечего почитать, и вы были стол благоразумны сохранив парочку статей та такой случай, смело вводите команду /rnd и наслаждайтесь чтивом.
В противном случае, передайте боту на хранения ссылку. Он любез но примет её и будет хранить в памяти (или в сердечке).

Очень важная информация: 
	При получении на хранения, бот не несёт ответственности за сданные ему ссылки. 
	После команды /rnd ссылка передаётся вам и забывается.`

const msgStart = "Привет! Ты попал в камеру хранения (как на вокзале).\n Она позволит тебе не таскаться по жизни с тяжёлым багажём ссылок на статьи.\n\n" + msgHelp

const (
	msgUnknownCommand = "Прошу прощения, я вас не понимаю."
	msgNoSavePage     = "У вас нет сохранённых ссылок."
	msgSaveCommand    = "Взял на хранение."
	msgAlreadyExists  = "Эта ссылка уже в хранилище."
)
