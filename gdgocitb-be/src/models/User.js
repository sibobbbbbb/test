import mongoose from 'mongoose';
import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';
import { config } from '../config/index.js';

const UserSchema = new mongoose.Schema({
  name: {
    type: String,
    required: [true, 'Name is required'],
    trim: true
  },
  email: {
    type: String,
    required: [true, 'Email is required'],
    unique: true,
    trim: true,
    lowercase: true,
    match: [
      /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/,
      'Please add a valid email'
    ]
  },
  access: {
    type: String,
    enum: ['Member', 'Buddy', 'Curriculum Admin', 'Professional Development Admin', 'Technical Admin'],
    required: [true, 'Access level is required']
  },
  paths: [{
    pathId: {
      type: mongoose.Schema.Types.ObjectId,
      ref: 'Path'
    },
    progress: [{
      moduleId: {
        type: mongoose.Schema.Types.ObjectId,
        ref: 'Module'
      },
      completed: {
        type: Boolean,
        default: false
      },
      // Problem sets completed for this module
      problemSetsCompleted: [{
        problemSetId: {
          type: mongoose.Schema.Types.ObjectId,
          ref: 'ProblemSet'
        },
        grade: {
          type: Number,
          default: 0
        },
        submissionUrl: String,
        submittedAt: Date
      }]
    }]
  }],
  certificates: [{
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Certificate'
  }],
  events: [{
    eventId: {
      type: mongoose.Schema.Types.ObjectId,
      refPath: 'events.eventType'
    },
    eventType: {
      type: String,
      enum: ['ProfessionalEvent', 'CommunityEvent']
    },
    rsvp: {
      type: Boolean,
      default: false
    },
    attended: {
      type: Boolean,
      default: false
    }
  }],
  createdAt: {
    type: Date,
    default: Date.now
  }
}, {
  timestamps: true
});

// Generate JWT Token
UserSchema.methods.getSignedJwtToken = function() {
  return jwt.sign(
    { id: this._id, access: this.access },
    config.jwt.secret,
    { expiresIn: config.jwt.expiresIn }
  );
};

const User = mongoose.models.User || mongoose.model('User', UserSchema);

export default User;