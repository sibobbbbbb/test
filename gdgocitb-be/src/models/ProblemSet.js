import mongoose from 'mongoose';

const ProblemSetSchema = new mongoose.Schema({
  pathId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Path',
    required: [true, 'Path ID is required']
  },
  problemSetTitle: {
    type: String,
    required: [true, 'Problem set title is required'],
    trim: true
  },
  description: {
    type: String,
    trim: true
  },
  video: {
    type: String,
    trim: true
  },
  submissionType: {
    type: String,
    enum: ['File', 'Link', 'Image', 'GOCI'],
    required: [true, 'Submission type is required']
  },
  accessLevel: {
    type: String,
    enum: ['Member', 'Buddy'],
    default: 'Member'
  },
  deadline: {
    type: Date
  },
  maxGrade: {
    type: Number,
    default: 100
  },
  passingGrade: {
    type: Number,
    default: 60
  },
  isManualGrading: {
    type: Boolean,
    default: false
  },
  order: {
    type: Number,
    default: 0
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

const ProblemSet = mongoose.model('ProblemSet', ProblemSetSchema);

export default ProblemSet;