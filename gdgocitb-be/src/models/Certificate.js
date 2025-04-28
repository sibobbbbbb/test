import mongoose from 'mongoose';

const CertificateSchema = new mongoose.Schema({
  userId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'User',
    required: [true, 'User ID is required']
  },
  pathId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Path',
    required: [true, 'Path ID is required']
  },
  certificateType: {
    type: String,
    enum: ['Member', 'Buddy'],
    required: [true, 'Certificate type is required']
  },
  name: {
    type: String,
    required: [true, 'Certificate name is required'],
    trim: true
  },
  issueDate: {
    type: Date,
    default: Date.now
  },
  certificateUrl: {
    type: String,
    trim: true
  },
  certificateId: {
    type: String,
    unique: true,
    required: [true, 'Certificate ID is required']
  },
  createdAt: {
    type: Date,
    default: Date.now
  },
  updatedAt: {
    type: Date,
    default: Date.now
  }
}, {
  timestamps: true
});

// Generate unique certificate ID
CertificateSchema.pre('save', function(next) {
  if (!this.certificateId) {
    this.certificateId = 'GDGOC-ITB-' + 
      this.certificateType.substring(0, 1).toUpperCase() + 
      '-' + 
      Math.floor(100000 + Math.random() * 900000) + 
      '-' + 
      new Date().getFullYear();
  }
  next();
});

const Certificate = mongoose.model('Certificate', CertificateSchema);

export default Certificate;