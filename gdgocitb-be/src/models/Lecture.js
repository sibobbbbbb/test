import mongoose from 'mongoose';

const LectureSchema = new mongoose.Schema({
  pathId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Path',
    required: [true, 'Path ID is required']
  },
  moduleId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Module',
    required: [true, 'Module ID is required']
  },
  title: {
    type: String,
    required: [true, 'Lecture title is required'],
    trim: true
  },
  notes: {
    type: String,
    trim: true
  },
  slides: {
    type: String,
    trim: true
  },
  sourceCode: {
    type: String,
    trim: true
  },
  order: {
    type: Number,
    default: 0
  },
  accessLevel: {
    type: String,
    enum: ['Member', 'Buddy'],
    default: 'Member'
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

const Lecture = mongoose.model('Lecture', LectureSchema);

export default Lecture;