import mongoose from 'mongoose';

const PathSchema = new mongoose.Schema({
  pathName: {
    type: String,
    required: [true, 'Path name is required'],
    trim: true
  },
  description: {
    type: String,
    trim: true
  },
  modules: [{
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Module'
  }],
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

const Path = mongoose.model('Path', PathSchema);

export default Path;