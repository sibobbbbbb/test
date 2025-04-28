import mongoose from 'mongoose';

const ModuleSchema = new mongoose.Schema({
  pathId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Path',
    required: [true, 'Path ID is required']
  },
  moduleName: {
    type: String,
    required: [true, 'Module name is required'],
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
  lectures: [{
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Lecture'
  }],
  problemSet: [{
    type: mongoose.Schema.Types.ObjectId,
    ref: 'ProblemSet'
  }],
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

const Module = mongoose.model('Module', ModuleSchema);

export default Module;