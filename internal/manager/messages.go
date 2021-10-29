package manager

type messageType int

const (
	messageTypeUnknown messageType = iota
	messageTypeRegistration
	messageTypeLectures
	messageTypeLecture
	messageTypeHomework
	messageTypeAllHomeworks
	messageTypeHomeworkPass
	messageTypeHelp
)

const (
	cmdRegistration = "регистрация"
	cmdLectures     = "лекции"
	cmdLecture      = "лекция"
	cmdHomework     = "получить дз"
	cmdAllHomeworks = "все дз"
	cmdHomeworkPass = "сдать дз"
	cmdReview       = "отзыв"
)

var (
	msgIDontKnowU           = "Я вас не знаю \xF0\x9F\x98\x94 Вы можете зарегистрироваться, отправив мне команду сообщение вида 'регистрация <your-email>'"
	msgWrongEmailFormat     = "Кажется, это неправильный email \xF0\x9F\x98\x94 Пожалуйста, пришлите email в формате example@edu.hse.ru"
	msgWelcomeUsername      = "Привет, %s! Вы успешно зарегистрированы \xF0\x9F\x8E\x89"
	msgLecturesList         = "Вот список лекций. Чтобы получить файл лекции, напишите 'лекция <номер лекции>', например, 'лекция 1'"
	msgCantGetLectureNumber = "Я не понял, какой номер лекции нужен \xF0\x9F\x98\x94 Чтобы получить файл лекции, напишите 'лекция <номер лекции>', например, 'лекция 1'"
	msgNoActiveHomework     = "Сейчас нет активных домашних заданий \xF0\x9F\x8E\x89"
	msgFileNotSent          = "Не могу найти файл \xF0\x9F\x98\x94 Чтобы сдать домашку, нужно отправить мне файл с тетрадкой."
	msgSmthWrong            = "Что-то пошло не так \xF0\x9F\x98\x94 Попробуйте еще раз"
	msgHomeworkPassed       = "Домашка успешно принята! \xF0\x9F\x98\x9C"
	msgHWExpired            = "Я всего лишь машина, так что ничего личного, но ДЗ просрочено, сдать его уже не получится \xF0\x9F\x98\x9E"
	msgDescribeCommands     = "Я бот курса дата-журналистики по программированию \xE2\x9C\x8C" +
		"С помощью меня вы можете получать лекции и домашние задания. Еще я умею принимать ваши тетрадки с домашкой. \n" +
		"\n\xE2\x9C\x85	Чтобы получить список лекций, напишите мне сообщение 'лекции'" +
		"\n\xE2\x9C\x85	Чтобы получить файл с конкретной лекцией, напишите мне сообщение 'лекция <номер лекции>, например, 'лекция 1'" +
		"\n\xE2\x9C\x85	Чтобы получить актуальное ДЗ, напишите мне сообщение 'получить дз', в ответ я отправлю файл с ним" +
		"\n\xE2\x9C\x85	Чтобы получить все ДЗ, напишите мне сообщение 'все дз'" +
		"\n\xE2\x9C\x85	Чтобы сдать ДЗ, отправьте мне файл с заданием" +
		"\n\xE2\x9C\x8F Чтобы оставить отзыв, напишите мне слово 'отзыв' и текст вашего отзыва в этом же сообщении"
	msgThankU = "Спасибо!"
)
