import express from 'express';
import { authenticateToken } from'./src/core/modules/auth-module/auth/auth.mjs';
import cors from 'cors';
import { UserController } from './src/core/modules/user-module/user-controller.mjs';
import { LoginController } from './src/core/modules/auth-module/login/login-controller.mjs';
import { sequelize } from './src/config/db.mjs';
import './src/core/models/UserModel.mjs';
import { EmployeeController } from './src/core/modules/employee-module/employee-controller.mjs';
import cookieParser from 'cookie-parser';
import fileUpload from 'express-fileupload';
import { ReportController } from './src/core/modules/report-module/report-controller.mjs';
import { EmailController } from './src/core/modules/email-module/email-controller.mjs';
import { SignatureController } from './src/core/modules/signature-module/signature-controller.mjs';
import fs from 'fs';
import path from 'path';
import crypto from 'crypto';
import { _log } from './src/common/helper/logger.mjs';
import helmet from 'helmet';

const ENCRYPTION_KEY = process.env.ENCRYPTION_KEY || 'zqufz8izbNPG4xwkrBl9f5kPtHLFrmhw';
const IV_LENGTH = process.env.IV_LENGTH || 16;
class Server {

	constructor() {
		this.app = express();
		this.config();
		this.generateKeyPair();
		this.routes();
		this.errorHandler();
		this.app.use(this.interceptRequest);
		this.app.use(helmet());
		this.app.use((req, res, next) => {
			res.header('Access-Control-Allow-Origin', 'http://localhost:5173');
			res.header('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept, Authorization');
			res.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
			res.header('Access-Control-Allow-Credentials', 'true');
			next();
		});
	}

	config() {
		const corsOptions = {
			origin: ['http://localhost:5173', 'http://172.21.0.7:5173'],
			credentials: true
		}
		this.app.use(cors(corsOptions));
		this.app.use(express.json());
		this.app.use(cookieParser());
		this.app.use(fileUpload({
			logger: console.log
		}));
	}

	syncTables() {
		sequelize.sync({ force: true }).then(() => {
			console.log('Database & tables created!');
		}).catch(error => {
			console.error('Unable to connect to the database:', error);
		});
	}

	routes() {
		this.userRoutes();
		this.loginRoutes();
		this.employeeRoutes();
		this.pdfRoutes();
		this.emailRoutes();
		this.signatureRoutes();
	}

	userRoutes() {
		const userController = UserController.getInstance();

		this.app.get('/api/users', authenticateToken, userController.getAllUsers);
		this.app.get('/api/users/:id', authenticateToken, userController.getUserById);
		this.app.post('/api/users', userController.createUser);
		this.app.put('/api/users/:id', authenticateToken, userController.updateUser);
		this.app.delete('/api/users/:id', authenticateToken, userController.deleteUser);
	}

	loginRoutes() {
		const loginController = LoginController.getInstance();

		this.app.post('/api/login', loginController.login);
		this.app.get('/api/check-token', authenticateToken, loginController.checkToken);
	}

	employeeRoutes() {
		const employeeController = EmployeeController.getInstance();

		this.app.post('/api/employees', authenticateToken, employeeController.createEmployee);
		this.app.get('/api/employees', authenticateToken, employeeController.getAllEmployees);
		this.app.get('/api/employees/:id', authenticateToken, employeeController.getOneEmployee);
		this.app.put('/api/employees/:id', authenticateToken, employeeController.updateEmployee);
		this.app.delete('/api/employees/:id', authenticateToken, employeeController.deleteEmployee);
	}

	pdfRoutes() {
		const reportController = ReportController.getInstance();

		this.app.post('/api/reports', authenticateToken, reportController.createReport);
		this.app.get('/api/reports', authenticateToken, reportController.getAllReports);
		this.app.get('/api/reports/signed', authenticateToken, reportController.getAllReportsSigned);
		this.app.get('/api/reports/:id', authenticateToken, reportController.getOneReport);
		this.app.delete('/api/reports/:id', authenticateToken, reportController.deleteReport);
	}

	emailRoutes() {
		const emailController = EmailController.getInstance();

		this.app.post('/api/email', authenticateToken, emailController.sendEmail);
	}

	signatureRoutes() {
		const signatureController = SignatureController.getInstance(this.privateKey, this.publicKey);

        this.app.post('/api/sign-document', authenticateToken, signatureController.signDocument);
        this.app.post('/api/verify-signature', authenticateToken, signatureController.verifySignature);
    }

	generateKeyPair() {
        const keysDirectory = path.resolve(process.env.KEYS_DIRECTORY);
        const privateKeyPath = path.join(keysDirectory, process.env.PRIVATE_KEY_PATH);
        const publicKeyPath = path.join(keysDirectory, process.env.PUBLIC_KEY_PATH);

        // Cria o diretório se ele não existir
        if (!fs.existsSync(keysDirectory)) {
            fs.mkdirSync(keysDirectory, { recursive: true });
        }

        if (!fs.existsSync(privateKeyPath) || !fs.existsSync(publicKeyPath)) {
            const { publicKey, privateKey } = crypto.generateKeyPairSync('rsa', {
                modulusLength: 2048,
            });

            fs.writeFileSync(privateKeyPath, this.encrypt(privateKey.export({ type: 'pkcs1', format: 'pem' })));
            fs.writeFileSync(publicKeyPath, this.encrypt(publicKey.export({ type: 'pkcs1', format: 'pem' })));
        }

        this.privateKey = this.decrypt(fs.readFileSync(privateKeyPath, 'utf8'));
        this.publicKey = this.decrypt(fs.readFileSync(publicKeyPath, 'utf8'));
    }
	
	//testes
	encrypt(text) {
        let iv = crypto.randomBytes(IV_LENGTH);
        let cipher = crypto.createCipheriv('aes-256-cbc', Buffer.from(ENCRYPTION_KEY), iv);
        let encrypted = cipher.update(text);
        encrypted = Buffer.concat([encrypted, cipher.final()]);
        return iv.toString('hex') + ':' + encrypted.toString('hex');
    }

	//testes 2
    decrypt(text) {
        let textParts = text.split(':');
        let iv = Buffer.from(textParts.shift(), 'hex');
        let encryptedText = Buffer.from(textParts.join(':'), 'hex');
        let decipher = crypto.createDecipheriv('aes-256-cbc', Buffer.from(ENCRYPTION_KEY), iv);
        let decrypted = decipher.update(encryptedText);
        decrypted = Buffer.concat([decrypted, decipher.final()]);
        return decrypted.toString();
    }

	errorHandler() {
		this.app.use((err, req, res, next) => {
		  console.error(err.stack);
		  res.status(500).send('Server Error');
		});
	}

	start() {
		this.app.listen(process.env.PORT || 3000, () => {
			_log('Server started on port 3000');
		});
	}

	interceptRequest(req, res, next) {
		_log(`Chegou uma requisição ${req.method} para a URL: ${req.url}`);
		next();
	}
}

const server = new Server();
server.start();