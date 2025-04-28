import mongoose from 'mongoose';

const ProfessionalEventSchema = new mongoose.Schema({
  title: {
    type: String,
    required: [true, 'Event title is required'],
    trim: true
  },
  subtitle: {
    type: String,
    trim: true
  },
  category: {
    type: String,
    enum: ['Google Developer Groups', 'Scholarship', 'Internship', 'Exchange'],
    required: [true, 'Category is required']
  },
  description: {
    type: String,
    trim: true
  },
  location: {
    type: String,
    trim: true
  },
  date: {
    type: Date,
    required: [true, 'Event date is required']
  },
  time: {
    start: {
      type: String,
      required: [true, 'Event start time is required']
    },
    end: {
      type: String,
      required: [true, 'Event end time is required']
    }
  },
  accessLevel: {
    type: String,
    enum: ['Member', 'Member and Buddy'],
    default: 'Member'
  },
  capacity: {
    type: Number
  },
  attendees: [{
    userId: {
      type: mongoose.Schema.Types.ObjectId,
      ref: 'User'
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
  },
  updatedAt: {
    type: Date,
    default: Date.now
  }
}, {
  timestamps: true
});

const ProfessionalEvent = mongoose.model('ProfessionalEvent', ProfessionalEventSchema);

export default ProfessionalEvent;