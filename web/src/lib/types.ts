// Domain types matching the Go backend

export interface School {
  id: number;
  name: string;
  subtitle: string;
  description: string;
  registeredAt: string;
  settings: SchoolSettings;
}

export interface SchoolSettings {
  color: string;
  domains: string[];
  contactInfo: ContactInfo;
  pages: Pages;
  showPaymentImages: boolean;
  logo: string;
  googleAnalyticsCode: string;
  fondy: Fondy;
  sendpulse: SendPulse;
  disableRegistration: boolean;
}

export interface ContactInfo {
  businessName: string;
  registrationNumber: string;
  address: string;
  email: string;
  phone: string;
}

export interface Pages {
  confidential: string;
  serviceAgreement: string;
  newsletterConsent: string;
}

export interface Fondy {
  merchantId: string;
  merchantPassword: string;
  connected: boolean;
}

export interface SendPulse {
  id: string;
  secret: string;
  listId: string;
  connected: boolean;
}

export interface Admin {
  id: number;
  name: string;
  email: string;
  schoolId: number;
}

export interface Course {
  id: number;
  name: string;
  description: string;
  imageUrl: string;
  color: string;
  published: boolean;
  modules?: Module[];
}

export interface Module {
  id: number;
  name: string;
  position: number;
  published: boolean;
  lessons?: Lesson[];
  survey?: Survey;
}

export interface Lesson {
  id: number;
  name: string;
  position: number;
  published: boolean;
  content?: LessonContent;
}

export interface LessonContent {
  id: number;
  lessonId: number;
  content: string;
}

export interface Package {
  id: number;
  name: string;
  description: string;
  benefits: string[];
  price: Price;
  modules: number[];
}

export interface Price {
  value: number;
  currency: string;
}

export interface Offer {
  id: number;
  name: string;
  description: string;
  price: Price;
  benefits: string[];
  paymentMethod: PaymentMethod;
  moduleId: number;
}

export interface PaymentMethod {
  usesProvider: boolean;
}

export interface Student {
  id: number;
  name: string;
  email: string;
  schoolId: number;
  blocked: boolean;
  offers?: Offer[];
}

export interface Order {
  id: number;
  studentId: number;
  offerId: number;
  status: string;
  amount: number;
  currency: string;
  createdAt: string;
  offer?: Offer;
  promo?: PromoInfo;
}

export interface PromoInfo {
  code: string;
  discount: number;
}

export interface PromoCode {
  id: number;
  code: string;
  discount: number;
  expiresAt: string;
  active: boolean;
  offerIds: number[];
}

export interface Survey {
  title: string;
  required: boolean;
  questions: SurveyQuestion[];
}

export interface SurveyQuestion {
  id: number;
  question: string;
  answerType: string;
  answerOptions: string[];
}

export interface SurveyResult {
  id: number;
  surveyId: number;
  studentId: number;
  answers: SurveyAnswer[];
  student?: { name: string; email: string };
}

export interface SurveyAnswer {
  questionId: number;
  answer: string;
}

export interface TokenResponse {
  accessToken: string;
  refreshToken: string;
}

export interface DataResponse<T> {
  data: T[];
  count: number;
}

export interface IdResponse {
  id: number;
}
